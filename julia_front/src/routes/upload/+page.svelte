<script>
	import { onMount } from 'svelte';
	import { isAuthenticated, fetchWithAuth } from '$lib/auth';
	import { goto } from '$app/navigation';

	// 사용자 정보 관련 상태
	let currentUser = null;
	let userRole = 'student';
	let userLoading = true;

	// 파일 업로드 관련 상태
	let selectedFile = null;
	let isUploading = false;
	let uploadProgress = 0;
	let uploadResult = null;
	let uploadError = null;

	// 파일 처리 관련 상태
	let isProcessing = false;
	let processingResult = null;
	let processingError = null;

	// 환경변수에서 파일 처리 서버 URL 가져오기
	const fileProcessorUrl = import.meta.env.VITE_FILE_PROCESSOR_URL || 'http://localhost:5001';

	onMount(async () => {
		// 로그인 상태 확인
		const authenticated = await isAuthenticated();
		if (!authenticated) {
			await goto('/');
			return;
		}
		
		await loadCurrentUser();
	});

	// 현재 사용자 정보 로드
	async function loadCurrentUser() {
		try {
			const res = await fetchWithAuth('/api/auth/me');
			if (res.ok) {
				currentUser = await res.json();
				userRole = currentUser.role || 'student';
				console.log('사용자 정보 로드 성공:', currentUser);
			} else {
				console.error('사용자 정보 로드 실패');
				userRole = 'student';
			}
		} catch (e) {
			console.error('사용자 정보 로드 오류:', e);
			userRole = 'student';
		} finally {
			userLoading = false;
		}
	}

	// 파일 선택 핸들러
	function handleFileSelect(event) {
		const file = event.target.files[0];
		if (file) {
			// 파일 확장자 검증
			const allowedExtensions = ['.hwp', '.hwpx'];
			const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
			
			if (!allowedExtensions.includes(fileExtension)) {
				alert('HWP 또는 HWPX 파일만 업로드할 수 있습니다.');
				event.target.value = '';
				return;
			}

			// 파일 크기 검증 (16MB)
			if (file.size > 16 * 1024 * 1024) {
				alert('파일 크기는 16MB를 초과할 수 없습니다.');
				event.target.value = '';
				return;
			}

			selectedFile = file;
			uploadResult = null;
			uploadError = null;
			processingResult = null;
			processingError = null;
		}
	}

	// 파일 업로드
	async function uploadFile() {
		if (!selectedFile) {
			alert('파일을 선택해주세요.');
			return;
		}

		isUploading = true;
		uploadProgress = 0;
		uploadError = null;

		try {
			const formData = new FormData();
			formData.append('file', selectedFile);

			// 파일 처리 서버로 직접 전송
			const response = await fetch(`${fileProcessorUrl}/api/process-docx`, {
				method: 'POST',
				body: formData
			});

			if (response.ok) {
				uploadResult = await response.json();
				uploadProgress = 100;
				showToast('파일이 성공적으로 처리되었습니다!', 'success');
			} else {
				const errorData = await response.json();
				uploadError = errorData.error || '파일 처리 중 오류가 발생했습니다.';
				showToast(uploadError, 'error');
			}
		} catch (error) {
			console.error('파일 업로드 오류:', error);
			uploadError = '파일 업로드 중 오류가 발생했습니다.';
			showToast(uploadError, 'error');
		} finally {
			isUploading = false;
		}
	}

	// 파일 다운로드
	async function downloadFile(filename) {
		try {
			const response = await fetch(`${fileProcessorUrl}/api/download/${filename}`);
			if (response.ok) {
				const blob = await response.blob();
				const url = window.URL.createObjectURL(blob);
				const a = document.createElement('a');
				a.href = url;
				a.download = filename;
				document.body.appendChild(a);
				a.click();
				window.URL.revokeObjectURL(url);
				document.body.removeChild(a);
				showToast('파일 다운로드가 시작되었습니다.', 'success');
			} else {
				showToast('파일 다운로드에 실패했습니다.', 'error');
			}
		} catch (error) {
			console.error('파일 다운로드 오류:', error);
			showToast('파일 다운로드 중 오류가 발생했습니다.', 'error');
		}
	}

	// 파일 초기화
	function resetUpload() {
		selectedFile = null;
		uploadResult = null;
		uploadError = null;
		processingResult = null;
		processingError = null;
		uploadProgress = 0;
		document.getElementById('fileInput').value = '';
	}

	// 로그아웃 핸들러
	async function handleLogout() {
		try {
			const res = await fetchWithAuth('/api/auth/logout', {
				method: 'POST',
			});

			if (res.ok) {
				localStorage.removeItem('authToken');
				goto('/');
			} else {
				console.error('로그아웃 실패');
				localStorage.removeItem('authToken');
				goto('/');
			}
		} catch (error) {
			console.error('로그아웃 오류:', error);
			localStorage.removeItem('authToken');
			goto('/');
		}
	}

	// 토스트 알림 함수
	function showToast(message, type = 'info') {
		const toast = document.createElement('div');
		toast.className = `toast toast-${type}`;
		toast.textContent = message;
		toast.style.cssText = `
			position: fixed;
			top: 20px;
			right: 20px;
			background: ${type === 'success' ? '#10b981' : type === 'error' ? '#ef4444' : '#3b82f6'};
			color: white;
			padding: 12px 20px;
			border-radius: 8px;
			z-index: 10000;
			font-size: 14px;
			font-weight: 500;
			box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
			animation: slideInRight 0.3s ease-out;
		`;
		
		document.body.appendChild(toast);
		
		setTimeout(() => {
			toast.style.animation = 'slideOutRight 0.3s ease-in forwards';
			setTimeout(() => {
				document.body.removeChild(toast);
			}, 300);
		}, 3000);
	}
</script>

<svelte:head>
	<title>파일 업로드 - Julia</title>
	<meta name="description" content="Julia 시스템 파일 업로드" />
</svelte:head>

<div class="upload-container">
	<header class="upload-header">
		<div class="header-left">
			<h1>파일 업로드</h1>
			<div class="user-info">
				<span>안녕하세요, {userLoading ? '로딩 중...' : currentUser?.id || '사용자'}님!</span>
				<span class="user-role {userRole}">{userRole === 'admin' ? '관리자' : '학생'}</span>
			</div>
		</div>
		<div class="header-right">
			<button class="btn ghost" on:click={() => goto('/dashboard')}>
				<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z" stroke="currentColor" fill="none" stroke-width="2"/>
					<path d="M8 5V3a2 2 0 012-2h4a2 2 0 012 2v2" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
				대시보드
			</button>
			<button class="btn ghost" on:click={() => goto('/calendar')}>
				<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
				캘린더
			</button>
			<button class="logout-button" on:click={handleLogout}>
				로그아웃
			</button>
		</div>
	</header>

	<main class="upload-content">
		{#if userRole !== 'admin'}
			<div class="access-denied">
				<h2>접근 권한이 없습니다</h2>
				<p>관리자만 파일 업로드 기능을 사용할 수 있습니다.</p>
				<button class="btn primary" on:click={() => goto('/dashboard')}>대시보드로 돌아가기</button>
			</div>
		{:else}
			<div class="upload-section">
				<div class="upload-card">
					<div class="card-header">
						<h2>📄 HWP/HWPX 파일 업로드</h2>
						<p>HWP 또는 HWPX 파일을 업로드하여 처리할 수 있습니다.</p>
					</div>

					<div class="upload-area">
						{#if !selectedFile}
							<div class="drop-zone" on:click={() => document.getElementById('fileInput').click()}>
								<svg width="48" height="48" viewBox="0 0 24 24" fill="none">
									<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2"/>
									<polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2"/>
									<line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2"/>
									<line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2"/>
									<polyline points="10,9 9,9 8,9" stroke="currentColor" stroke-width="2"/>
								</svg>
								<h3>파일을 선택하거나 여기에 드래그하세요</h3>
								<p>HWP, HWPX 파일만 지원됩니다 (최대 16MB)</p>
							</div>
						{:else}
							<div class="file-selected">
								<div class="file-info">
									<svg width="24" height="24" viewBox="0 0 24 24" fill="none">
										<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2"/>
										<polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2"/>
									</svg>
									<div class="file-details">
										<h4>{selectedFile.name}</h4>
										<p>{(selectedFile.size / 1024 / 1024).toFixed(2)} MB</p>
									</div>
								</div>
								<button class="btn ghost small" on:click={resetUpload}>
									<svg width="16" height="16" viewBox="0 0 24 24">
										<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2"/>
									</svg>
									제거
								</button>
							</div>
						{/if}

						<input
							id="fileInput"
							type="file"
							accept=".hwp,.hwpx"
							on:change={handleFileSelect}
							style="display: none;"
						/>
					</div>

					{#if isUploading}
						<div class="upload-progress">
							<div class="progress-bar">
								<div class="progress-fill" style="width: {uploadProgress}%"></div>
							</div>
							<p>파일을 처리하고 있습니다...</p>
						</div>
					{/if}

					{#if uploadError}
						<div class="error-message">
							<svg width="20" height="20" viewBox="0 0 24 24">
								<path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="currentColor" stroke-width="2"/>
							</svg>
							<p>{uploadError}</p>
						</div>
					{/if}

					{#if uploadResult}
						<div class="success-message">
							<svg width="20" height="20" viewBox="0 0 24 24">
								<path d="M9 12l2 2 4-4M21 12c0 4.97-4.03 9-9 9s-9-4.03-9-9 4.03-9 9-9 9 4.03 9 9z" stroke="currentColor" stroke-width="2"/>
							</svg>
							<p>{uploadResult.message}</p>
							<button class="btn primary" on:click={() => downloadFile(uploadResult.filename)}>
								<svg width="16" height="16" viewBox="0 0 24 24">
									<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3" stroke="currentColor" stroke-width="2"/>
								</svg>
								파일 다운로드
							</button>
						</div>
					{/if}

					<div class="upload-actions">
						<button 
							class="btn primary" 
							on:click={uploadFile}
							disabled={!selectedFile || isUploading}
						>
							{#if isUploading}
								<svg width="16" height="16" viewBox="0 0 24 24" class="spinning">
									<path d="M21 12a9 9 0 11-6.219-8.56" stroke="currentColor" stroke-width="2"/>
								</svg>
								처리 중...
							{:else}
								<svg width="16" height="16" viewBox="0 0 24 24">
									<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2"/>
									<polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2"/>
								</svg>
								파일 업로드 및 처리
							{/if}
						</button>
						{#if selectedFile}
							<button class="btn ghost" on:click={resetUpload} disabled={isUploading}>
								초기화
							</button>
						{/if}
					</div>
				</div>

				<div class="info-card">
					<h3>📋 사용 안내</h3>
					<ul>
						<li>HWP, HWPX 파일만 업로드할 수 있습니다.</li>
						<li>파일 크기는 최대 16MB까지 지원됩니다.</li>
						<li>업로드된 파일은 Windows 서버에서 처리됩니다.</li>
						<li>처리 완료 후 다운로드 링크가 제공됩니다.</li>
					</ul>
				</div>
			</div>
		{/if}
	</main>
</div>

<style>
	:root {
		--bg: #f6f7fb;
		--card: #ffffff;
		--text: #1f2937;
		--muted: #6b7280;
		--line: #e5e7eb;
		--brand: #3b82f6;
		--brand-600: #2563eb;
		--success: #10b981;
		--error: #ef4444;
		--warning: #f59e0b;
		--danger: #ef4444;
	}

	* {
		box-sizing: border-box;
	}

	.upload-container {
		min-height: 100vh;
		background: var(--bg);
		color: var(--text);
	}

	.upload-header {
		background: var(--card);
		border-bottom: 1px solid var(--line);
		padding: 1rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
	}

	.header-left h1 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		font-weight: 700;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.user-role {
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-weight: 600;
		font-size: 0.75rem;
	}

	.user-role.admin {
		background: var(--brand);
		color: white;
	}

	.user-role.student {
		background: var(--muted);
		color: white;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		border-radius: 0.5rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		text-decoration: none;
	}

	.btn:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	.btn.ghost {
		background: transparent;
	}

	.btn.primary {
		background: var(--brand);
		border-color: var(--brand);
		color: white;
	}

	.btn.small {
		padding: 0.25rem 0.5rem;
		font-size: 0.875rem;
	}

	.logout-button {
		padding: 0.5rem 1rem;
		background: var(--danger);
		color: white;
		border: none;
		border-radius: 0.375rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.logout-button:hover {
		background: #dc2626;
	}

	.upload-content {
		padding: 2rem;
		max-width: 800px;
		margin: 0 auto;
	}

	.access-denied {
		text-align: center;
		padding: 3rem 2rem;
		background: var(--card);
		border-radius: 1rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.access-denied h2 {
		color: var(--error);
		margin-bottom: 1rem;
	}

	.upload-section {
		display: grid;
		gap: 2rem;
	}

	.upload-card, .info-card {
		background: var(--card);
		border-radius: 1rem;
		padding: 2rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.card-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.card-header h2 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		font-weight: 700;
	}

	.card-header p {
		color: var(--muted);
		margin: 0;
	}

	.upload-area {
		margin-bottom: 2rem;
	}

	.drop-zone {
		border: 2px dashed var(--line);
		border-radius: 1rem;
		padding: 3rem 2rem;
		text-align: center;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.drop-zone:hover {
		border-color: var(--brand);
		background: color-mix(in srgb, var(--brand) 5%, transparent);
	}

	.drop-zone svg {
		color: var(--muted);
		margin-bottom: 1rem;
	}

	.drop-zone h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.drop-zone p {
		color: var(--muted);
		margin: 0;
	}

	.file-selected {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background: color-mix(in srgb, var(--brand) 10%, transparent);
		border: 1px solid var(--brand);
		border-radius: 0.5rem;
	}

	.file-info {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.file-info svg {
		color: var(--brand);
	}

	.file-details h4 {
		margin: 0 0 0.25rem 0;
		font-weight: 600;
	}

	.file-details p {
		margin: 0;
		color: var(--muted);
		font-size: 0.875rem;
	}

	.upload-progress {
		margin-bottom: 1rem;
	}

	.progress-bar {
		width: 100%;
		height: 0.5rem;
		background: var(--line);
		border-radius: 0.25rem;
		overflow: hidden;
		margin-bottom: 0.5rem;
	}

	.progress-fill {
		height: 100%;
		background: var(--brand);
		transition: width 0.3s ease;
	}

	.upload-progress p {
		text-align: center;
		color: var(--muted);
		margin: 0;
	}

	.error-message, .success-message {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: 0.5rem;
		margin-bottom: 1rem;
	}

	.error-message {
		background: color-mix(in srgb, var(--error) 10%, transparent);
		border: 1px solid var(--error);
		color: var(--error);
	}

	.success-message {
		background: color-mix(in srgb, var(--success) 10%, transparent);
		border: 1px solid var(--success);
		color: var(--success);
	}

	.error-message svg, .success-message svg {
		flex-shrink: 0;
	}

	.error-message p, .success-message p {
		margin: 0;
		font-weight: 600;
	}

	.upload-actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
	}

	.info-card h3 {
		margin: 0 0 1rem 0;
		font-size: 1.125rem;
		font-weight: 700;
	}

	.info-card ul {
		margin: 0;
		padding-left: 1.5rem;
	}

	.info-card li {
		margin-bottom: 0.5rem;
		color: var(--muted);
	}

	.spinning {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	@keyframes slideInRight {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes slideOutRight {
		from {
			transform: translateX(0);
			opacity: 1;
		}
		to {
			transform: translateX(100%);
			opacity: 0;
		}
	}

	@media (max-width: 768px) {
		.upload-header {
			flex-direction: column;
			gap: 1rem;
			text-align: center;
		}

		.header-right {
			flex-wrap: wrap;
			justify-content: center;
		}

		.upload-content {
			padding: 1rem;
		}

		.upload-card, .info-card {
			padding: 1.5rem;
		}

		.drop-zone {
			padding: 2rem 1rem;
		}

		.upload-actions {
			flex-direction: column;
		}
	}
</style>
