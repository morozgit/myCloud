import axios from 'axios';

const apiBase = import.meta.env.VITE_API_BASE;

export async function uploadFile(formData) {
  console.log('formData', formData);
  const response = await axios.post(`${apiBase}/files/upload`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
}
