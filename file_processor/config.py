import os

class Config:
    WINDOWS_SERVER_URL = os.getenv('WINDOWS_SERVER_URL', 'http://192.168.45.102:5000')
    UPLOAD_FOLDER = '/app/uploads'
    ALLOWED_EXTENSIONS = {'hwp', 'hwpx'}
    MAX_CONTENT_LENGTH = 16 * 1024 * 1024  # 16MB
    PROCESS_ENDPOINT = "/api/synonym"
