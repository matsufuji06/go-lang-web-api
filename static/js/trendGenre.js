document.addEventListener('click', (e) => {
  const country = e.target.closest('.svgMap-country');
  if (!country) return;

  const iso = country.dataset.id;
  console.log(iso);
  fetch(`http://localhost:8000/api/v1/analytics/genres?country=${iso}`)
    .then(res => {
      if (res.status === 403) {
        throw new Error("この国ではYouTubeは利用できません");
      } else if (res.status === 500) {
        throw new Error("動画の取得に失敗しました");
      } else if (res.status === 429) {
        throw new Error("APIの呼び出し上限に達しました");
      }
      console.log('status:', res.status);
      console.log('response:', res);
      return res.json();
    })
    .then(data => {
      renderPie(data.genres);
      renderVideos(data.genres);
    })
    .catch(err => {
      alert("動画の取得に失敗しました");
      console.error(err)
    });
})

function renderPie(genres) {
  const tbody = document.querySelector('#genre-pie tbody');
  tbody.innerHTML = '';
  const legend = document.getElementById('genre-legend');
  legend.innerHTML = '';

  let total = 0;

  for (let i = 0; i < genres.length; i++) {
    total = total + genres[i].ratio;
  }

  let start = 0;

  genres.forEach(genre => {
    const value = genre.ratio / total;
    const end = start + value;
    // 円グラフ
    tbody.insertAdjacentHTML(
      'beforeend',
      `
      <tr>
        <td style="--start:${start}; --end:${end};">
          <span class="data">${(value * 100).toFixed(1)}%</span>
        </td>
      </tr>
      `
    );
    // 凡例
    legend.insertAdjacentHTML(
      'beforeend',
      `<li>${genre.genre}：${(value * 100).toFixed(1)}%</li>`
    );
    start = end;
  });
}

function renderVideos(genres) {
  const container = document.getElementById('genre-videos');
  container.innerHTML = '';

  genres.forEach(genre => {
    // ジャンル見出し
    container.insertAdjacentHTML(
      'beforeend',
      `<h5 class="mb-3">${genre.genre}</h5>`
    );

    const ul = document.createElement('ul');
    ul.innerHTML = "";

    genre.topVideos.forEach(video => {
      ul.insertAdjacentHTML(
        'beforeend',
        `
        <li>
          <a href="${video.url}" target="_blank">
            ${video.url}
          </a>
        </li>
        `
      );
    });

    container.appendChild(ul);
  });
}