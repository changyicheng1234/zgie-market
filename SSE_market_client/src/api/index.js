import request from '@/utils/request';

// 1获取文件夹列表
function getFileListApi(data) {
  let formData = new FormData();
  formData.append("path", data);
  return request.post('/auth/getFileList',formData);
}

// 2新增文件夹
function addFolderApi(data) {
  let formData = new FormData();
  formData.append("path", data);
  return request.post('/auth/addFolder',formData);
}

// 3上传文件
function uploadFile(query){
  const config = {
    timeout:60000
  }
  return request.post('/auth/uploadFile',query,config);
}

// 4获取文件url
function getUrl(filename){
  let formData = new FormData();
  formData.append("filename", filename);
  return request.post('/auth/getObjectUrl',formData);
}
// 6删除文件
function delFileApi(data) {
  let formData = new FormData();
  formData.append("path", data);
  return request.post('/auth/fileDelete',formData);
}
// 7删除文件夹
function  delFolderApi(data) {
  let formData = new FormData();
  formData.append("path", data);
  return request.post('/auth/folderDelete',formData);
}
// 8下载数据
function downloadDataApi(data) {
  return request('./data/download.text');
}

// 9文件夹移动
function fileMoveApi(data) {
  return request('./data/move-file.json');
}

export {
  getFileListApi, // 1获取文件夹列表
  addFolderApi, // 2新增文件夹
  uploadFile, // 获取全部文件夹
  getUrl, // 4获取文件url
  delFileApi, // 6删除文件
  delFolderApi,//7 删除文件夹
  downloadDataApi, // 8下载数据
  fileMoveApi, // 9文件夹移动
};
