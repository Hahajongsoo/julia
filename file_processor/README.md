# File Processor Service

Julia 프로젝트의 파일 처리 서비스입니다. HWP/HWPX 파일을 Windows 서버로 전송하여 처리하는 기능을 제공합니다.

## 기능

- HWP/HWPX 파일 업로드 및 처리
- Windows 서버로 파일 전송
- 처리된 파일 다운로드
- 비동기 파일 처리
- 동의어 파일 업로드

## API 엔드포인트

### 파일 처리
- `POST /api/process-docx` - 파일 업로드 및 즉시 처리
- `POST /api/process-docx-async` - 비동기 파일 처리
- `GET /api/download/<filename>` - 처리된 파일 다운로드

### 기타
- `GET /health` - 헬스체크
- `POST /api/upload-synonyms` - 동의어 파일 업로드

## 환경 변수

- `WINDOWS_SERVER_URL`: Windows 서버 URL (기본값: http://192.168.45.102:5000)
- `FLASK_ENV`: Flask 환경 (production/development)
- `FLASK_DEBUG`: 디버그 모드 (True/False)

## Docker 실행

### Development
```bash
docker-compose -f docker-compose.dev.yml up file-processor-dev
```

### Production
```bash
docker-compose -f docker-compose.prod.yml up file-processor
```

## 포트

- **5001**: 파일 처리 서비스 포트
- Production에서는 외부에서 접근 가능하도록 포트가 열려있습니다.

## Windows 서버 연결

Production 환경에서 Windows 서버(192.168.45.102:5000)에 연결할 수 있도록 포트 5001이 외부에 노출되어 있습니다.
