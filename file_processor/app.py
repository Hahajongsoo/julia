from flask import Flask, request, jsonify, send_file
from flask_cors import CORS
import os
import uuid
import logging
from werkzeug.utils import secure_filename
import requests
from datetime import datetime

app = Flask(__name__)
CORS(app)

# 로깅 설정
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 설정
UPLOAD_FOLDER = '/app/uploads'
ALLOWED_EXTENSIONS = {'hwp', 'hwpx'}
WINDOWS_SERVER_URL = os.getenv('WINDOWS_SERVER_URL', 'http://192.168.45.102:5000')
PROCESS_ENDPOINT = "/api/synonym"

# 업로드 폴더 생성
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

def send_to_windows_server(file_path):
    """
    Windows 서버로 hwp 또는 hwpx 파일을 전송하고 결과를 받아옵니다.
    """
    try:
        # 파일 확장자 체크
        ext = os.path.splitext(file_path)[1].lower()
        if ext not in ['.hwp', '.hwpx']:
            raise Exception('hwp 또는 hwpx 파일만 전송할 수 있습니다.')
        
        # Windows 서버의 URL과 엔드포인트 설정
        if not WINDOWS_SERVER_URL:
            raise Exception('WINDOWS_SERVER_URL 환경변수가 설정되지 않았습니다.')
        
        windows_server_url = f"{WINDOWS_SERVER_URL.rstrip('/')}{PROCESS_ENDPOINT}"
        logger.info(f"Sending file to Windows server at: {windows_server_url}")
        
        # 파일 전송
        with open(file_path, 'rb') as f:
            files = {'file': f}
            response = requests.post(windows_server_url, files=files, timeout=300)
        
        if response.status_code != 200:
            error_msg = response.text if response.text else '알 수 없는 오류'
            if 'application/json' in response.headers.get('Content-Type', ''):
                try:
                    error_data = response.json()
                    error_msg = error_data.get('error', error_msg)
                except:
                    pass
            raise Exception(f'Windows 서버 오류: {error_msg}')
        
        # 결과 파일 저장
        unique_id = uuid.uuid4().hex[:8]
        original_name = os.path.basename(file_path)
        filename = f"{os.path.splitext(original_name)[0]}_{unique_id}.hwp"
        output_path = os.path.join(UPLOAD_FOLDER, filename)
        
        with open(output_path, 'wb') as f:
            f.write(response.content)
        
        return {
            'message': '파일이 성공적으로 처리되었습니다.',
            'file_path': output_path,
            'filename': filename,
            'download_url': f'/api/download/{filename}'
        }
        
    except requests.RequestException as e:
        raise Exception(f'통신 오류: {str(e)}')
    except Exception as e:
        raise Exception(f'처리 중 오류 발생: {str(e)}')

@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({'status': 'healthy', 'timestamp': datetime.now().isoformat()})

@app.route('/api/process-docx', methods=['POST'])
def process_docx():
    """파일 업로드 및 처리"""
    try:
        if 'file' not in request.files:
            return jsonify({'error': '파일이 전송되지 않았습니다.'}), 400
        
        file = request.files['file']
        if file.filename == '':
            return jsonify({'error': '선택된 파일이 없습니다.'}), 400
        
        if not allowed_file(file.filename):
            return jsonify({'error': '지원하지 않는 파일 형식입니다. .hwp, .hwpx 파일만 허용됩니다.'}), 400
        
        # 파일을 임시 저장
        filename = secure_filename(file.filename)
        filepath = os.path.join(UPLOAD_FOLDER, filename)
        file.save(filepath)
        
        # Windows 서버로 전송 및 처리
        result = send_to_windows_server(filepath)
        
        # 임시 파일 삭제
        try:
            os.remove(filepath)
        except:
            pass
        
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Error processing file: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/api/process-docx-async', methods=['POST'])
def process_docx_async():
    """비동기 파일 처리"""
    try:
        data = request.get_json()
        filename = data.get('filename')
        
        if not filename:
            return jsonify({'error': '파일명이 전달되지 않았습니다.'}), 400
        
        filepath = os.path.join(UPLOAD_FOLDER, filename)
        if not os.path.exists(filepath):
            return jsonify({'error': '업로드된 파일을 찾을 수 없습니다.'}), 400
        
        # Windows 서버로 전송 및 처리
        result = send_to_windows_server(filepath)
        
        # 임시 파일 삭제
        try:
            os.remove(filepath)
        except:
            pass
        
        return jsonify({
            'message': result['message'],
            'download_url': result['download_url']
        })
        
    except Exception as e:
        logger.error(f"Error processing file async: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/api/download/<filename>')
def download_processed_file(filename):
    """처리된 파일 다운로드"""
    try:
        file_path = os.path.join(UPLOAD_FOLDER, filename)
        if not os.path.exists(file_path):
            return jsonify({'error': '파일을 찾을 수 없습니다.'}), 404
        
        return send_file(
            file_path,
            as_attachment=True,
            download_name=filename
        )
    except Exception as e:
        logger.error(f"Error downloading file: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/api/upload-synonyms', methods=['POST'])
def upload_synonyms():
    """동의어 파일 업로드"""
    try:
        if 'file' not in request.files:
            return jsonify({'error': 'No file part'}), 400
        
        file = request.files['file']
        if file.filename == '':
            return jsonify({'error': 'No selected file'}), 400
        
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            file_path = os.path.join(UPLOAD_FOLDER, filename)
            file.save(file_path)
            
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    data = f.read()
                
                return jsonify({
                    'message': 'File uploaded and processed successfully',
                    'data': data
                })
            except Exception as e:
                return jsonify({'error': str(e)}), 500
            finally:
                # 임시 파일 삭제
                if os.path.exists(file_path):
                    os.remove(file_path)
        
        return jsonify({'error': 'Invalid file type'}), 400
        
    except Exception as e:
        logger.error(f"Error uploading synonyms: {str(e)}")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001, debug=False)
