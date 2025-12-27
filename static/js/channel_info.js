const config = {
    baseUrl: "http://localhost:8000/api/v1/channels?query=",
    parentId: "channel-info-container",
    buttonId: "channel-info-search-button"
}

const button = document.getElementById(config.buttonId);
button.addEventListener("click", function(e) {
    e.preventDefault(); // フォームのデフォルト送信を防止
    const query = document.getElementById("channel-info-query").value;
    fetchChannelInfo(query);
})

// クエリに基づいてチャンネル情報を取得する関数
function fetchChannelInfo(query) {
    const url = config.baseUrl + encodeURIComponent(query);
    fetch(url)
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    alert(err.message || "Error occurred");
                    throw new Error("API error");
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("API result:", data);
            displayChannelInfo(data);
        })
        .catch(error => {
            console.error("Error fetching channel info:", error);
        });
}

function displayChannelInfo(data) {
    const container = document.getElementById(config.parentId);
    container.innerHTML = ""; // 既存の内容をクリア

    const card = document.createElement("div");
    card.className = "card";
    card.className = "card shadow-sm mb-3";
    card.innerHTML = `
        <div class="card-body">
            <div class="d-flex align-items-start gap-3">
            
            <!-- サムネイル -->
            <img 
                src="${data.thumbnail}" 
                alt="Channel Thumbnail"
                class="rounded"
                style="width: 96px; height: 96px; object-fit: cover;"
            />

            <!-- テキスト -->
            <div>
                <h5 class="card-title mb-2">${data.title}</h5>
                <p class="mb-1 text-muted">登録者数: ${formatNumberJP(data.subscriberCount)}</p>
                <p class="mb-1 text-muted">総再生数: ${formatNumberJP(data.viewCount)}</p>
                <p class="mb-1 text-muted">動画本数: ${formatNumberJP(data.videoCount)}</p>
                <p class="mb-0 text-muted">平均再生数: ${formatNumberJP(data.averageViews)}</p>
            </div>

            </div>
        </div>
        `;

    container.appendChild(card);
}

// 大きい数字を見やすく表示するための関数 (例: 1200000 -> "120万")
function formatNumberJP(num) {
    if (num >= 100000000) {
        return (num / 100000000).toFixed(1).replace(/\.0$/, "") + "億";
    }
    if (num >= 10000) {
        return (num / 10000).toFixed(1).replace(/\.0$/, "") + "万";
    }
    return num.toLocaleString();
}