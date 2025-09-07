<script>
	import { createEventDispatcher } from 'svelte';
	import { fetchWithAuth } from '$lib/auth';
	import { API_ENDPOINTS } from '$lib/config';

	const dispatch = createEventDispatcher();

	export let isVisible = false;
	export let userRole = 'student';

	let memoContent = '';
	let memoLastSaved = null;
	let autoSaveTimeout = null;

	// 메모 패널 열기/닫기
	function togglePanel() {
		isVisible = !isVisible;
		if (isVisible) {
			loadMemo();
		} else {
			// 자동 저장 취소
			if (autoSaveTimeout) {
				clearTimeout(autoSaveTimeout);
			}
		}
		dispatch('toggle', { isVisible });
	}

	function closePanel() {
		isVisible = false;
		// 자동 저장 취소
		if (autoSaveTimeout) {
			clearTimeout(autoSaveTimeout);
		}
		dispatch('close');
	}

	// 메모 로드
	async function loadMemo() {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO);
			if (response.ok) {
				const data = await response.json();
				memoContent = data.content || '';
				memoLastSaved = data.updated_at || null;
			} else {
				memoContent = '';
				memoLastSaved = null;
			}
		} catch (error) {
			console.error('메모 로드 오류:', error);
			memoContent = '';
			memoLastSaved = null;
		}
	}

	// 메모 내용 변경 핸들러 (자동 저장)
	function handleMemoChange() {
		// 기존 자동 저장 타이머 취소
		if (autoSaveTimeout) {
			clearTimeout(autoSaveTimeout);
		}
		
		// 3초 후 자동 저장
		autoSaveTimeout = setTimeout(() => {
			if (memoContent.trim()) {
				saveMemo(true); // 자동 저장 플래그
			}
		}, 3000);
	}

	// 메모 저장
	async function saveMemo(isAutoSave = false) {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ content: memoContent })
			});
			
			if (response.ok) {
				memoLastSaved = new Date().toISOString();
				if (!isAutoSave) {
					// 수동 저장일 때만 알림
					dispatch('toast', { message: '메모가 저장되었습니다.', type: 'success' });
				}
			}
		} catch (error) {
			console.error('메모 저장 오류:', error);
			if (!isAutoSave) {
				dispatch('toast', { message: '메모 저장에 실패했습니다.', type: 'error' });
			}
		}
	}

	// 메모 전체 삭제
	async function clearMemo() {
		if (!confirm('정말로 메모를 모두 삭제하시겠습니까?')) return;
		
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO, {
				method: 'DELETE'
			});
			
			if (response.ok) {
				memoContent = '';
				memoLastSaved = null;
				dispatch('toast', { message: '메모가 삭제되었습니다.', type: 'success' });
			}
		} catch (error) {
			console.error('메모 삭제 오류:', error);
			dispatch('toast', { message: '메모 삭제에 실패했습니다.', type: 'error' });
		}
	}

	// 마지막 저장 시간 포맷팅
	function formatLastSaved(dateString) {
		const date = new Date(dateString);
		const now = new Date();
		const diff = now - date;
		
		if (diff < 60000) { // 1분 미만
			return '방금 전';
		} else if (diff < 3600000) { // 1시간 미만
			return `${Math.floor(diff / 60000)}분 전`;
		} else if (diff < 86400000) { // 1일 미만
			return `${Math.floor(diff / 3600000)}시간 전`;
		} else {
			return date.toLocaleDateString('ko-KR', {
				month: 'short',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		}
	}

	// 외부에서 호출할 수 있도록 함수 노출
	export { togglePanel, closePanel };
</script>

{#if isVisible}
	<div class="memo-panel">
		<div class="memo-header">
			<h3>📝 관리자 메모장</h3>
			<button class="memo-close-btn" on:click={closePanel} aria-label="메모장 닫기">
				<svg width="20" height="20" viewBox="0 0 24 24">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
			</button>
		</div>
		
		<div class="memo-content">
			<!-- 메모 편집 영역 -->
			<div class="memo-editor">
				<textarea 
					bind:value={memoContent}
					placeholder="여기에 메모를 자유롭게 작성하세요..."
					class="memo-textarea"
					on:input={handleMemoChange}
				></textarea>
			</div>
			
			<!-- 메모 액션 버튼들 -->
			<div class="memo-actions">
				<button 
					class="btn primary memo-save-btn" 
					on:click={() => saveMemo(false)}
					disabled={!memoContent.trim()}
				>
					<svg width="16" height="16" viewBox="0 0 24 24">
						<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" stroke="currentColor" fill="none" stroke-width="2"/>
						<polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" fill="none" stroke-width="2"/>
						<polyline points="7,3 7,8 15,8" stroke="currentColor" fill="none" stroke-width="2"/>
					</svg>
					저장
				</button>
				
				<button 
					class="btn ghost memo-clear-btn" 
					on:click={clearMemo}
					disabled={!memoContent.trim()}
				>
					<svg width="16" height="16" viewBox="0 0 24 24">
						<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0v14M10 11v6M14 11v6" stroke="currentColor" fill="none" stroke-width="2"/>
					</svg>
					전체 삭제
				</button>
			</div>
			
			<!-- 메모 정보 -->
			<div class="memo-info">
				<div class="memo-stats">
					<span class="memo-char-count">{memoContent.length}자</span>
					{#if memoLastSaved}
						<span class="memo-last-saved">마지막 저장: {formatLastSaved(memoLastSaved)}</span>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	:root {
		--bg: #f6f7fb;
		--card: #ffffff;
		--text: #1f2937;
		--muted: #6b7280;
		--line: #e5e7eb;
		--brand: #3b82f6;
		--sun: #ef4444;
	}

	.btn {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		border-radius: 10px;
		padding: 8px 12px;
		font-weight: 600;
		cursor: pointer;
		transition: 0.15s ease;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.btn:hover {
		transform: translateY(-1px);
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

	.btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* 메모 패널 스타일 */
	.memo-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: 400px;
		height: 100vh;
		background: var(--card);
		border-left: 1px solid var(--line);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		box-shadow: -4px 0 20px rgba(0, 0, 0, 0.1);
		transform: translateX(100%);
		animation: slideInRight 0.3s ease-out forwards;
	}

	@keyframes slideInRight {
		to {
			transform: translateX(0);
		}
	}

	.memo-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		border-bottom: 1px solid var(--line);
		background: var(--bg);
	}

	.memo-header h3 {
		margin: 0;
		font-size: 18px;
		font-weight: 700;
		color: var(--text);
	}

	.memo-close-btn {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 8px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.memo-close-btn:hover {
		background: var(--line);
		color: var(--text);
	}

	.memo-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 20px;
	}

	.memo-editor {
		flex: 1;
		display: flex;
		flex-direction: column;
		margin-bottom: 20px;
	}

	.memo-textarea {
		flex: 1;
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 16px;
		font-family: inherit;
		font-size: 14px;
		line-height: 1.6;
		resize: none;
		background: var(--card);
		color: var(--text);
		transition: all 0.2s ease;
	}

	.memo-textarea:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.memo-textarea::placeholder {
		color: var(--muted);
		opacity: 0.7;
	}

	.memo-actions {
		display: flex;
		gap: 12px;
		margin-bottom: 16px;
	}

	.memo-save-btn,
	.memo-clear-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		flex: 1;
		justify-content: center;
	}

	.memo-clear-btn {
		color: var(--sun);
		border-color: var(--sun);
	}

	.memo-clear-btn:hover {
		background: var(--sun);
		color: white;
	}

	.memo-info {
		border-top: 1px solid var(--line);
		padding-top: 16px;
	}

	.memo-stats {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 12px;
		color: var(--muted);
	}

	.memo-char-count {
		font-weight: 600;
	}

	.memo-last-saved {
		opacity: 0.8;
	}

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.memo-panel {
			width: 100%;
			left: 0;
			right: 0;
		}
		
		.memo-actions {
			flex-direction: column;
		}
		
		.memo-save-btn,
		.memo-clear-btn {
			flex: none;
		}
	}
</style>
