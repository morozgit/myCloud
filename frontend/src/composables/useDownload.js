import { ref } from 'vue';
import toastr from 'toastr';
import axios from 'axios';

export const useDownload = () => {
  const isDownloading = ref(false);

  const downloadFile = async (item) => {
    isDownloading.value = true;
    try {
      const payload = {
        path: item.path,
        name: item.name
      };

      await axios.post('/api/files/download', payload);
      toastr.success('Файл добавлен в очередь на скачивание.');
    } catch (error) {
      console.error('Ошибка при отправке запроса на скачивание:', error);
      toastr.error('Ошибка: файл не был добавлен в очередь.');
    } finally {
      isDownloading.value = false;
    }
  };

  return { downloadFile, isDownloading };
};
