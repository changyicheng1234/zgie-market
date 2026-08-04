// 定义手机号规则
const telephoneValidator = (value) => {
  if (value === '') {
    return true; // 空白手机号视为合法
  }
  return /^1[2|3|4|5|6|7|8|9]\d{9}$/.test(value);
};
// 定义邮箱规则
const emailValidator = (value) => /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/.test(value);

export default {
  telephoneValidator,
  emailValidator,
};
