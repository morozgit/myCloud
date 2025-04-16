// src/composables/useDownload.js
import axios from 'axios';

/**
 * Отправка задания на скачивание файла через RabbitMQ
 * @param {Object} item - объект файла (содержит путь и имя)
 */
export const useDownload = () => {
  const downloadFile = async (item) => {
    console.log("TYTA")
    try {
      const payload = {
        path: item.path,    // полный путь до файла
        name: item.name     // имя файла
      };
      console.log("payload", payload)
      await axios.post('/api/navigation/download', payload);
      alert('Файл добавлен в очередь на скачивание.');
    } catch (error) {
      console.error('Ошибка при отправке запроса на скачивание:', error);
      alert('Не удалось отправить файл на скачивание.');
    }
  };

  return { downloadFile };
};
