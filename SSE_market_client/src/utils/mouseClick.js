export default {
  /* eslint-disable no-param-reassign */
  install(Vue) {
    Vue.prototype.$mouseClick = function () {
      let clickCount = 0;
      const content = ['中山大学', '❤SYSU❤', '软件工程学院', '❤SSE❤', '❤博学❤', '❤审问❤', '❤慎思❤', '❤明辨❤', '❤笃行❤'];
      const colors = ['red'];
      const htmlElement = document.getElementsByTagName('html')[0];
      const bodyElement = document.getElementsByTagName('body')[0];

      function getRandomInt(max) {
        return Math.floor(Math.random() * Math.floor(max));
      }

      htmlElement.onclick = function (e) {
        const element = document.createElement('b');
        element.style.color = colors[getRandomInt(colors.length)];
        element.style.zIndex = '9999';
        element.style.position = 'absolute';
        element.style.select = 'none';
        element.style.userSelect = 'none';
        const x = e.pageX;
        const y = e.pageY;
        element.style.left = `${x - 10}px`;
        element.style.top = `${y - 20}px`;

        if (clickCount >= content.length) {
          clickCount = 0;
        }

        element.innerText = content[clickCount];
        clickCount += 1;
        element.style.fontSize = '18px'; // 字体大小统一
        let increase = 0;

        setTimeout(() => {
          const anim = setInterval(() => {
            if (increase === 150) {
              clearInterval(anim);
              bodyElement.removeChild(element);
            }
            element.style.top = `${y - 20 - increase}px`;
            element.style.opacity = (150 - increase) / 120;
            increase += 1;
          }, 8);
        }, 70);

        bodyElement.appendChild(element);
      };
    };
  },
};
