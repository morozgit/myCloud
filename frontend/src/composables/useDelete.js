import { ref } from 'vue';
import toastr from 'toastr';
import axios from 'axios';

export const useDelete = () => {
  const isDownloading = ref(false);

  const deleteFile = async (item) => {
    isDownloading.value = true;
    try {
      const payload = {
        path: item.path,
        name: item.name
      };

      const response = await axios.post('/mycloud/api/files/delete', payload);
      console.log(response);
      toastr.success('Файл удалён');
    } catch (error) {
      console.error('Ошибка delete', error);
      toastr.error('Ошибка при удалении файла');
    } finally {
      isDownloading.value = false;
    }
  };

  return { deleteFile }; // <-- ЭТО ВАЖНО
};
