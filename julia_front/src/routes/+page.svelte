<script>
	import { goto } from '$app/navigation';
	import { API_ENDPOINTS } from '$lib/config';

  
	let id = '';
	let password = '';
	let remember = false;
	let showPwd = false;
	let loading = false;
	let errorMsg = '';	  

	async function onSubmit(e) {
	  e.preventDefault();
	  errorMsg = '';
	  loading = true;
	  try {
		const res = await fetch(API_ENDPOINTS.LOGIN, {
		  method: 'POST',
		  headers: { 'Content-Type': 'application/json' },
		  credentials: 'include', // 세션/쿠키 기반 인증 시 유지
		  body: JSON.stringify({ id, password, remember })
		});
  
		if (!res.ok) {
		  const t = await res.text().catch(()=>'');
		  throw new Error(t || '로그인에 실패했습니다.');
		}
  
		await goto('/calendar');
	  } catch (err) {
		errorMsg = (err && err.message) ? err.message : '알 수 없는 오류가 발생했습니다.';
	  } finally {
		loading = false;
	  }
	}
  </script>
  
  <svelte:head>
	<title>로그인 | Julia</title>
	<meta name="description" content="Julia 시스템 로그인" />
  </svelte:head>
  
  <main>
	<div class="form-container login-card">
	  <h1 class="form-title">로그인</h1>
  
	  {#if errorMsg}
		<div class="error-message" role="alert">{errorMsg}</div>
	  {/if}
  
	  <form class="login-form" on:submit={onSubmit} novalidate>
		<div class="form-group">
		  <label class="form-label" for="id">ID</label>
		  <input
			id="id"
			class="form-input"
			type="text"
			placeholder="이름"
			bind:value={id}
			required
			autocomplete="username"
			inputmode="text"
		  />
		</div>
  
		<div class="form-group">
		  <label class="form-label" for="password">비밀번호</label>
  
		  <div class="pw-wrap">
			{#if showPwd}
			  <!-- type이 고정된 input (text) -->
			  <input
				id="password"
				class="form-input pw-input"
				type="text"
				placeholder="비밀번호"
				bind:value={password}
				required
				autocomplete="current-password"
			  />
			{:else}
			  <!-- type이 고정된 input (password) -->
			  <input
				id="password"
				class="form-input pw-input"
				type="password"
				placeholder="••••••••"
				bind:value={password}
				required
				autocomplete="current-password"
			  />
			{/if}
  
			<button
			  type="button"
			  class="pw-eye"
			  on:click={() => (showPwd = !showPwd)}
			  aria-label={showPwd ? '비밀번호 숨기기' : '비밀번호 표시'}
			>
			  {#if showPwd}
				<!-- eye-off -->
				<svg viewBox="0 0 24 24" aria-hidden="true">
				  <path d="M3 3l18 18M10.6 10.65A2 2 0 0012 14a2 2 0 001.4-.58M9.88 7.7A8.6 8.6 0 0112 7c5.5 0 9 5 9 5a15.8 15.8 0 01-3.07 3.43M6.6 6.56A15.3 15.3 0 003 12s3.5 5 9 5a8.9 8.9 0 003.1-.55" fill="none" stroke="currentColor" stroke-width="1.8"/>
				</svg>
			  {:else}
				<!-- eye -->
				<svg viewBox="0 0 24 24" aria-hidden="true">
				  <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12zm10 3a3 3 0 110-6 3 3 0 010 6z" fill="none" stroke="currentColor" stroke-width="1.8"/>
				</svg>
			  {/if}
			</button>
		  </div>
		</div>
  

  
		<button class="form-button" type="submit" disabled={loading}>
		  {#if loading}
			<span class="spinner" aria-hidden="true"></span> 로그인 중…
		  {:else}
			로그인
		  {/if}
		</button>
	  </form>

	</div>
  </main>
  
  <style>
	/* 로그인 페이지 전용 스타일 */
	.login-card {
		backdrop-filter: saturate(150%) blur(4px);
		border: 1px solid rgba(0, 0, 0, 0.06);
	}
  </style>
  