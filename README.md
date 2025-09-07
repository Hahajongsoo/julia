# Julia 프로젝트 Docker Compose 설정

이 프로젝트는 프론트엔드(SvelteKit)와 백엔드(Go)를 Docker Compose로 관리합니다.

## 프로젝트 구조

```
julia/
├── julia_front/          # SvelteKit 프론트엔드
├── julia_back/           # Go 백엔드
├── docker-compose.yml    # 전체 서비스 (개발 + 프로덕션)
├── docker-compose.dev.yml    # 개발 환경만
├── docker-compose.prod.yml   # 프로덕션 환경만
└── README.md
```

## 사용법

### 1. 개발 환경 실행

```bash
# 개발 환경만 실행 (권장)
docker-compose -f docker-compose.dev.yml up -d

# 또는 전체 서비스 실행
docker-compose up julia-backend-dev julia-frontend-dev postgres -d
```

### 2. 프로덕션 환경 실행

```bash
# 방법 1: 프로덕션 전용 파일 사용
# 1. SSL 인증서 초기화 (최초 1회만)
# init-letsencrypt.sh 파일에서 이메일 주소를 수정한 후 실행
./init-letsencrypt.sh

# 2. 프로덕션 환경 실행
docker-compose -f docker-compose.prod.yml up -d

# 방법 2: 전체 통합 파일 사용
# 1. SSL 인증서 초기화 (최초 1회만)
./init-letsencrypt.sh

# 2. 프로덕션 서비스만 실행
docker-compose up julia-backend julia-frontend postgres nginx certbot -d

# 3. SSL 인증서 갱신 (필요시)
./renew-ssl.sh
```

### 3. 서비스 접속

- **프론트엔드 (개발)**: http://localhost:5173
- **백엔드 (개발)**: http://localhost:8081
- **프론트엔드 (프로덕션)**: https://hapcky.me
- **백엔드 API (프로덕션)**: https://hapcky.me/api
- **PostgreSQL**: 내부 네트워크만 (포트 충돌 방지)

### 4. 서비스 관리

```bash
# 서비스 상태 확인
docker-compose ps

# 로그 확인
docker-compose logs -f [서비스명]

# 서비스 중지
docker-compose down

# 볼륨까지 삭제
docker-compose down -v

# 특정 서비스만 재시작
docker-compose restart nginx

# SSL 인증서 갱신 후 nginx 재시작
./renew-ssl.sh
```

### 5. 개발 시 유용한 명령어

```bash
# 특정 서비스만 재시작
docker-compose restart julia-backend

# 특정 서비스만 재빌드
docker-compose build julia-frontend

# 실시간 로그 확인
docker-compose logs -f julia-backend julia-frontend
```

## 환경 변수

### PostgreSQL 설정
- **이미지**: postgres:latest
- **비밀번호**: 1234
- **타임존**: Asia/Seoul
- **포트**: 5432
- **볼륨**: 기존 PostgreSQL 컨테이너의 볼륨 재사용

### 백엔드 설정
- **개발 모드**: GIN_MODE=debug
- **프로덕션 모드**: GIN_MODE=release
- **포트**: 8080

### 프론트엔드 설정
- **개발 모드**: NODE_ENV=development (Vite dev server)
- **프로덕션 모드**: NODE_ENV=production (정적 파일 서버)
- **개발 포트**: 5173
- **프로덕션 포트**: 3000 (serve 정적 파일 서버)
- **API URL**: Docker 환경에서는 `/api`, 로컬에서는 `http://localhost:8080`

## 네트워크

모든 서비스는 `julia-network`라는 브리지 네트워크를 통해 통신합니다.

## 볼륨

- 기존 PostgreSQL 컨테이너의 볼륨을 외부 볼륨으로 재사용하여 데이터 영속성 보장
- 볼륨 ID: `94d218c76603151366cc7c91c857ce664cf06917481aaf59b5bbcb4d97d87be8`
- `external: true` 설정으로 기존 볼륨을 참조

## SSL 인증서

- Let's Encrypt를 통한 무료 SSL 인증서 자동 발급
- 도메인: `hapcky.me`, `www.hapcky.me`
- 인증서 자동 갱신 지원
- Nginx를 통한 HTTPS 리버스 프록시

## 문제 해결

### 포트 충돌
- **개발 환경**: 포트 8081, 5173, 5432 사용
- **프로덕션 환경**: nginx만 외부 포트 80, 443 사용 (내부 서비스는 expose만 사용)
- 포트 충돌이 발생하면 docker-compose.yml 파일에서 포트 매핑을 변경하세요.

### 권한 문제
Docker 권한이 없는 경우:
```bash
sudo usermod -aG docker $USER
# 로그아웃 후 다시 로그인
```

### 빌드 실패
캐시를 지우고 다시 빌드:
```bash
docker-compose build --no-cache
```
