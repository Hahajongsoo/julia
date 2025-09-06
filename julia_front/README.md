# Julia Frontend

Julia 시스템의 Svelte 기반 프론트엔드 애플리케이션입니다.

## 기능

- 로그인 페이지 (`/`)
- 대시보드 페이지 (`/dashboard`)
- Go 백엔드와의 API 연동

## 개발 환경 설정

### 로컬 개발

```bash
# 의존성 설치
npm install

# 개발 서버 실행
npm run dev
```

### Docker를 사용한 개발

```bash
# 개발 환경 실행
docker-compose up julia-frontend-dev

# 프로덕션 빌드 및 실행
docker-compose up julia-frontend
```

## API 연동

프론트엔드는 `http://localhost:8081`에서 실행되는 Go 백엔드와 연동됩니다.

### 주요 API 엔드포인트

- `POST /auth/login` - 로그인
- `POST /auth/logout` - 로그아웃

### 인증 시스템

- 쿠키 기반 세션 관리
- 자동 쿠키 포함 요청 (`credentials: 'include'`)
- 세션 만료 시 자동 로그아웃

## 환경 변수

프로젝트 루트에 `.env` 파일을 생성하고 다음 환경변수들을 설정하세요:

```bash
# VAPID Public Key (Push Notification용)
VITE_VAPID_PUBLIC=your_vapid_public_key_here

# Backend Host URL
VITE_BACKEND_HOST=http://localhost:8080

# 기타 환경변수들
VITE_APP_TITLE=Julia
VITE_APP_VERSION=1.0.0
VITE_DEBUG_MODE=false
```

- `VITE_VAPID_PUBLIC` - Push Notification을 위한 VAPID Public Key
- `VITE_BACKEND_HOST` - 백엔드 서버 호스트 URL
- `API_BASE_URL` - 백엔드 API 기본 URL (기본값: `/api`)

## 빌드

```bash
# 프로덕션 빌드
npm run build

# 빌드 결과 미리보기
npm run preview
```
