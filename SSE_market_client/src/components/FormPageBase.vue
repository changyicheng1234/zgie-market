<template>
  <div class='form-page-view'>
    <vue-particles
      class="particles-bg"
      color="#FFC0CB"
      :particle-opacity="0.7"
      :particles-number="100"
      shape-type="star"
      :particle-size="4"
      lines-color="#8DD1FE"
      :lines-width="1"
      :line-linked="true"
      :line-opacity="0.4"
      :lines-distance="150"
      :move-speed="3"
      :hover-effect="true"
      hover-mode="grab"
      :click-effect="true"
      click-mode="push"
    />
    <div class="d-flex flex-wrap justify-content-center align-items-center" style="height: 100vh;">
      <b-card class="form-card">
        <slot />
      </b-card>
    </div>
  </div>
</template>

<style>
.particles-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.form-page-view {
  background-image: url('https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/background.jpg');
  background-size: cover;
  background-position: center;
  position: fixed;
  top: 0;
  left: 0;
  margin-top: 0px !important;
  height: 100vh;
  width: 100vw;
}

.form-card {
  height: fit-content;
  max-height: 100vh;
  width: 600px;
  overflow: auto;
}
</style>
