// ================================
// 設定
// ================================
const API_BASE_URL = "http://localhost:8000/api/v1/videos";
const POLLING_INTERVAL = 60_000; // 60秒

// ================================
// 状態
// ================================
let lastPublishedAt = null;
let pollingTimer = null;

// 表示済み動画ID
const displayedVideoIds = new Set();

// ================================
// DOM取得
// ================================
const form = document.getElementById("textSearchForm");
const keywordInput = document.getElementById("keywordInput");
const limitSelect = document.getElementById("limitSelect");
const timeline = document.getElementById("videoTimeline");
const pollingStatus = document.getElementById("pollingStatus");

// ================================
// API通信
// ================================
// 動画データを取得
async function fetchVideos({ keyword, limit, since }) {
  // クエリパラメータを構築
  const params = new URLSearchParams({
    keyword, limit,
  });

  // sinceパラメータがあれば追加
  if (since) {
    params.append("since", since);
  }

  // API呼び出し
  const response = await fetch(`${API_BASE_URL}?${params.toString()}`);
  if (!response.ok) {
    throw new Error(`APIエラー: ${response.status}`);
  }

  return response.json();
}

// ================================
// UI生成
// ================================
function createVideoCard(video) {
  const newBadge = video.isNew
    ? `<span class="badge bg-success position-absolute top-0 start-0 m-2">NEW</span>`
    : "";

  return `
    <div class="col-md-6 col-lg-4">
      <div class="card h-100 shadow-sm position-relative">
        ${newBadge}
        <img
          src="${video.thumbnails.high.url}"
          class="card-img-top"
          alt="thumbnail"
        />
        <div class="card-body">
          <h5 class="card-title">${video.title}</h5>
          <p class="text-muted mb-1">${video.channelName}</p>
          <p class="small text-muted">
            公開日: ${new Date(video.publishedAt).toLocaleString()}
          </p>
          <a
            href="${video.url}"
            target="_blank"
            class="btn btn-sm btn-outline-danger"
          >
            YouTubeで見る
          </a>
        </div>
      </div>
    </div>
  `;
}

// ================================
// 描画
// ================================
// タイムラインに動画を描画
function renderTimeline(videos, append = false) {
  // 既存のコンテンツをクリア
  if (!append) {
    timeline.innerHTML = "";
    displayedVideoIds.clear();
  }

  // forEachで追加する順序を制御するため、あえて配列を反転する
  const videosReversed = append ? [...videos].reverse() : videos;

  // 動画をタイムラインに追加
  videosReversed.forEach(video => {
    displayedVideoIds.add(video.videoId);

    // appendがtrueなら先頭に追加、falseなら末尾に追加
    // ※true: ポーリング時、
    // ※false: 初回検索時
    timeline.insertAdjacentHTML(
      append ? "afterbegin" : "beforeend",
      createVideoCard(video)
    );
  });
}

// ================================
// ポーリング
// ================================
// ポーリング開始
function startPolling(keyword, limit) {
  // 既存のポーリングを停止
  stopPolling();

  pollingTimer = setInterval(async () => {
    // ポーリング状態を表示
    const startTime = Date.now();
    showPollingStatus();

    try {
      // 新しい動画を取得
      const data = await fetchVideos({ keyword, limit, since: lastPublishedAt });

      // 表示済み動画IDと照合して新しい動画のみ抽出し、isNewフラグを付与
      const newItems = data.items
        .filter(video => !displayedVideoIds.has(video.videoId))
        .map(video => ({
          ...video,
          isNew: true,
        }));

      // 新しい動画がなければ終了
      if (newItems.length == 0) {
        return
      }

      // 投稿日時で降順ソート
      const sortedItems = sortByPublishedAtDesc(newItems);

      // 新しい動画があればタイムラインに追加
      renderTimeline(sortedItems, true);

      // 最新の投稿日時を更新
      lastPublishedAt = sortedItems[0].publishedAt;

    } catch (error) {
      console.error("ポーリングエラー:", error);

    } finally {
      const elapsed = Date.now() - startTime;
      const remaining = Math.max(0, 3000 - elapsed);

      setTimeout(() => {
        hidePollingStatus();
      }, remaining);

    }
  }, POLLING_INTERVAL);
}

// ポーリング停止
function stopPolling() {
  if (pollingTimer) {
    clearInterval(pollingTimer);
    pollingTimer = null;
  }
}

// ================================
// 検索イベント
// ================================
form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const keyword = keywordInput.value.trim();
  const limit = limitSelect.value;

  if (keyword === "") {
    alert("キーワードを入力してください。");
    return;
  }

  // 既存のポーリングを停止
  stopPolling();

  // 検索ごとに前回投稿日時リセット
  lastPublishedAt = null;

  try {
    // 動画を取得して描画
    const data = await fetchVideos({ keyword, limit });

    // 投稿日時で降順ソートし、isNewフラグをfalseに設定（初回はisNew=false）
    const sortedItems = sortByPublishedAtDesc(data.items)
      .map(video => ({ ...video, isNew: false }));

    // タイムラインを更新
    renderTimeline(sortedItems);
    if (data.items.length > 0) {
      // 最新の投稿日時を更新
      lastPublishedAt = sortedItems[0].publishedAt;
    }
    // ポーリング開始
    startPolling(keyword, limit);

  } catch (error) {
    console.error("Search error:", error);
    alert("動画の取得に失敗しました。");
  }
});

// ================================
// ユーティリティ関数
// ================================
// 投稿日時で降順ソート
function sortByPublishedAtDesc(videos) {
  return [...videos].sort((a, b) => {
    return new Date(b.publishedAt) - new Date(a.publishedAt);
  });
}

// ポーリング状態を表示
function showPollingStatus() {
  pollingStatus.classList.remove("d-none");
}

// ポーリング状態を非表示
function hidePollingStatus() {
  pollingStatus.classList.add("d-none");
}
