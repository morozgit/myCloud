import axios from 'axios';

export async function uploadFile(formData) {
  console.log('formData', formData);
  const response = await axios.post('/mycloud/api/files/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
}
