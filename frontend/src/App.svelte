<script>
  import logo from './assets/images/logo-universal.png'
  import {UploadImages, StitchImages, DownloadResult} from '../wailsjs/go/main/App.js'

  let uploadedImages = []
  let stitching = false
  let stitchResult = null
  let stitchMode = 'horizontal'
  let errorMessage = ''

  function handleFileInput(event) {
    const files = event.target.files
    if (files.length > 0) {
      const images = []
      for (let i = 0; i < files.length; i++) {
        const reader = new FileReader()
        reader.onload = function(e) {
          images.push(e.target.result)
          if (images.length === files.length) {
            uploadImages(images)
          }
        }
        reader.readAsDataURL(files[i])
      }
    }
  }

  function uploadImages(images) {
    UploadImages(images).then(result => {
      if (result.success) {
        uploadedImages = result.images
        errorMessage = ''
      } else {
        errorMessage = result.error
      }
    }).catch(error => {
      errorMessage = '上传图片失败: ' + error.message
    })
  }

  function stitch() {
    // 检查不同拼接模式所需的图片数量
    if (stitchMode === 'grid-2x2' && uploadedImages.length !== 4) {
      errorMessage = '2x2网格模式需要4张图片'
      return
    } else if ((stitchMode === 'grid-3x1' || stitchMode === 'grid-1x3') && uploadedImages.length !== 3) {
      errorMessage = '3x1或1x3网格模式需要3张图片'
      return
    } else if (uploadedImages.length < 2) {
      errorMessage = '请至少上传2张图片'
      return
    }
    stitching = true
    StitchImages(uploadedImages, stitchMode).then(result => {
      stitching = false
      if (result.success) {
        stitchResult = result.image
        errorMessage = ''
      } else {
        errorMessage = result.error
      }
    }).catch(error => {
      stitching = false
      errorMessage = '拼接图片失败: ' + error.message
    })
  }

  function download() {
    if (stitchResult) {
      DownloadResult().then(result => {
        if (!result.success) {
          errorMessage = result.error
        }
      }).catch(error => {
        errorMessage = '下载结果失败: ' + error.message
      })
    }
  }

  function removeImage(index) {
    uploadedImages.splice(index, 1)
  }

  function clearAll() {
    uploadedImages = []
    stitchResult = null
    errorMessage = ''
  }
</script>

<main>
  <h1>图片拼接工具</h1>
  
  <div class="section upload-section">
    <h2>上传图片</h2>
    <label class="upload-btn" for="file-input">
      选择图片
      <input id="file-input" type="file" multiple accept="image/*" on:change={handleFileInput} />
    </label>
    {#if errorMessage}
      <div class="error-message">{errorMessage}</div>
    {/if}
  </div>

  <div class="section images-section">
    <h2>已上传图片</h2>
    <div class="settings" style="margin-bottom: 15px;">
      <button class="stitch-btn" on:click={clearAll} style="background-color: #f44336;">
        清除所有
      </button>
    </div>
    <div class="images-container">
      {#each uploadedImages as image, index}
        <div class="image-item">
          <img src={image} alt="已上传图片" />
          <button class="remove-btn" on:click={() => removeImage(index)}>×</button>
        </div>
      {/each}
      {#if uploadedImages.length === 0}
        <div class="no-images">还没有上传图片</div>
      {/if}
    </div>
  </div>

  <div class="section stitch-section">
    <h2>拼接设置</h2>
    <div class="settings">
      <label for="stitch-mode">拼接模式:</label>
      <select id="stitch-mode" bind:value={stitchMode}>
        <option value="horizontal">水平拼接</option>
        <option value="vertical">垂直拼接</option>
        <option value="grid-2x2">2x2网格</option>
        <option value="grid-3x1">3x1网格</option>
        <option value="grid-1x3">1x3网格</option>
      </select>
      <button class="stitch-btn" on:click={stitch} disabled={stitching || uploadedImages.length < 2}>
        {stitching ? '拼接中...' : '拼接图片'}
      </button>
    </div>
  </div>

  {#if stitchResult}
    <div class="section result-section">
      <h2>拼接结果</h2>
      <div class="result-container">
        <img src={stitchResult} alt="拼接后的图片" />
      </div>
      <button class="download-btn" on:click={download}>下载结果</button>
    </div>
  {/if}
</main>

<style>
  :root {
    --primary-color: #4CAF50;
    --primary-dark: #388E3C;
    --primary-light: #81C784;
    --secondary-color: #2196F3;
    --text-color: #333;
    --text-light: #757575;
    --background-color: #f5f5f5;
    --card-background: #ffffff;
    --border-radius: 8px;
    --box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    --transition: all 0.3s ease;
  }

  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background-color: var(--background-color);
    color: var(--text-color);
    line-height: 1.6;
  }

  main {
    max-width: 1000px;
    margin: 0 auto;
    padding: 20px;
  }

  h1 {
    text-align: center;
    color: var(--primary-dark);
    margin-bottom: 30px;
    font-size: 2.5rem;
    font-weight: 600;
  }

  h2 {
    color: var(--text-color);
    margin: 25px 0 15px;
    font-size: 1.5rem;
    font-weight: 500;
  }

  .section {
    background-color: var(--card-background);
    border-radius: var(--border-radius);
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: var(--box-shadow);
    transition: var(--transition);
  }

  .section:hover {
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
  }

  .upload-section {
    text-align: center;
  }

  input[type="file"] {
    display: none;
  }

  .upload-btn {
    display: inline-block;
    padding: 12px 24px;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: var(--transition);
    margin-bottom: 15px;
  }

  .upload-btn:hover {
    background-color: var(--primary-dark);
    transform: translateY(-2px);
  }

  .error-message {
    color: #f44336;
    margin-top: 10px;
    padding: 10px;
    background-color: rgba(244, 67, 54, 0.1);
    border-radius: var(--border-radius);
    border-left: 4px solid #f44336;
  }

  .images-container {
    display: flex;
    flex-wrap: wrap;
    gap: 15px;
    margin: 15px 0;
  }

  .image-item {
    position: relative;
    width: 160px;
    height: 160px;
    border: 2px solid #e0e0e0;
    border-radius: var(--border-radius);
    overflow: hidden;
    transition: var(--transition);
  }

  .image-item:hover {
    border-color: var(--primary-color);
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }

  .image-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .remove-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: rgba(255, 0, 0, 0.8);
    color: white;
    border: none;
    border-radius: 50%;
    width: 28px;
    height: 28px;
    cursor: pointer;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: var(--transition);
    opacity: 0;
  }

  .image-item:hover .remove-btn {
    opacity: 1;
  }

  .remove-btn:hover {
    background-color: rgba(255, 0, 0, 1);
    transform: scale(1.1);
  }

  .stitch-section .settings {
    display: flex;
    align-items: center;
    gap: 15px;
    margin: 15px 0;
    flex-wrap: wrap;
  }

  .stitch-section label {
    font-weight: 500;
    color: var(--text-light);
  }

  select {
    padding: 10px;
    border: 1px solid #e0e0e0;
    border-radius: var(--border-radius);
    font-size: 1rem;
    background-color: white;
    transition: var(--transition);
  }

  select:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
  }

  .stitch-btn, .download-btn {
    padding: 12px 24px;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--border-radius);
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: var(--transition);
  }

  .stitch-btn:hover, .download-btn:hover {
    background-color: var(--primary-dark);
    transform: translateY(-2px);
  }

  .stitch-btn:disabled {
    background-color: #e0e0e0;
    color: #9e9e9e;
    cursor: not-allowed;
    transform: none;
  }

  .result-section {
    text-align: center;
  }

  .result-container {
    margin: 20px 0;
    border: 1px solid #e0e0e0;
    border-radius: var(--border-radius);
    padding: 20px;
    background-color: white;
    box-shadow: var(--box-shadow);
  }

  .result-container img {
    max-width: 100%;
    height: auto;
    border-radius: var(--border-radius);
  }

  .download-btn {
    margin-top: 20px;
    background-color: var(--secondary-color);
  }

  .download-btn:hover {
    background-color: #1976D2;
  }

  .no-images {
    text-align: center;
    color: var(--text-light);
    padding: 40px 0;
    font-style: italic;
  }

  @media (max-width: 768px) {
    main {
      padding: 10px;
    }

    h1 {
      font-size: 2rem;
    }

    .stitch-section .settings {
      flex-direction: column;
      align-items: flex-start;
    }

    .stitch-btn, .download-btn {
      width: 100%;
      text-align: center;
    }

    .image-item {
      width: 120px;
      height: 120px;
    }
  }

</style>
