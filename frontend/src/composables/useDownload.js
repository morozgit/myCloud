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

      // Отправляем запрос на сервер для получения ссылки на скачивание
      const response = await axios.post('/api/files/download', payload);
      
      // Получаем ссылку для скачивания из ответа
      const downloadUrl = response.data.download_url;

      if (downloadUrl) {
        // Создаём временный элемент <a> для скачивания файла
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = item.name; // Устанавливаем имя файла для скачивания

        // Добавляем элемент на страницу, кликаем по нему и удаляем
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        toastr.success('Файл добавлен в очередь на скачивание.');
      } else {
        toastr.error('Не удалось получить ссылку для скачивания.');
      }
    } catch (error) {
      console.error('Ошибка при отправке запроса на скачивание:', error);
      toastr.error('Ошибка: файл не был добавлен в очередь.');
    } finally {
      isDownloading.value = false;
    }
  };

  return { downloadFile, isDownloading };
};
