import { useRef, useState } from 'react';
import { usePostUpload, usePostUploadJobIdComplete } from '@/shared/api/generated/imgpdf';
import pdfIcon from '@/shared/assets/pdf.svg';
import './PdfUploadPage.css';

export default function PdfUploadPage() {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const initUploadMutation = usePostUpload();
  const completeUploadMutation = usePostUploadJobIdComplete();

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const files = e.dataTransfer.files;
    if (files && files[0]) {
      setSelectedFile(files[0]);
      setUploadError(null);
      setUploadSuccess(false);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files[0]) {
      setSelectedFile(files[0]);
      setUploadError(null);
      setUploadSuccess(false);
    }
  };

  const handleRemoveFile = () => {
    setSelectedFile(null);
    setUploadError(null);
    setUploadSuccess(false);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleUploadClick = () => {
    if (selectedFile) {
      void handleUpload(selectedFile);
    }
  };

  const handleUpload = async (file: File) => {
    try {
      setIsUploading(true);
      setUploadError(null);

      // 1. Инициализируем загрузку
      const initResponse = await initUploadMutation.mutateAsync();

      if (initResponse.status !== 201) {
        throw new Error('Не удалось инициализировать загрузку');
      }

      const { job_id, upload_url } = initResponse.data;

      if (!job_id || !upload_url) {
        throw new Error('Неверный ответ сервера');
      }

      // 2. Загружаем файл по presigned URL
      const uploadResponse = await fetch(upload_url, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': 'application/pdf',
        },
      });

      if (!uploadResponse.ok) {
        throw new Error('Не удалось загрузить файл');
      }

      // 3. Завершаем загрузку
      const completeResponse = await completeUploadMutation.mutateAsync({ jobId: job_id });

      if (completeResponse.status !== 202) {
        throw new Error('Не удалось завершить загрузку');
      }

      setUploadSuccess(true);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Ошибка загрузки:', error);
      setUploadError(error instanceof Error ? error.message : 'Произошла ошибка при загрузке');
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="pdf-upload-page">
      <div className="upload-container">
        <h1 className="upload-title">Загрузка PDF</h1>

        <div
          className={`drop-zone ${isDragging ? 'dragging' : ''}`}
          onDragEnter={handleDragEnter}
          onDragLeave={handleDragLeave}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
        >
          <div className="drop-zone-content">
            <svg className="upload-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M7 18V20C7 21.1046 7.89543 22 9 22H15C16.1046 22 17 21.1046 17 20V18M12 2V16M12 2L8 6M12 2L16 6"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>

            <p className="drop-zone-text">{isDragging ? 'Отпустите файл здесь' : 'Перетащите PDF файл сюда'}</p>

            <p className="drop-zone-subtext">или</p>

            <label className="file-input-label">
              Выберите файл
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdf,application/pdf"
                onChange={handleFileSelect}
                className="file-input"
                disabled={isUploading}
              />
            </label>
          </div>
        </div>

        {selectedFile && (
          <div className="selected-file">
            <div className="file-info">
              <img className="file-icon" src={pdfIcon} alt="PDF" />

              <div className="file-details">
                <p className="file-name">{selectedFile.name}</p>
                <p className="file-size">{(selectedFile.size / 1024 / 1024).toFixed(2)} МБ</p>
              </div>
            </div>

            <button className="remove-button" onClick={handleRemoveFile} disabled={isUploading}>
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M18 6L6 18M6 6L18 18"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </button>
          </div>
        )}

        {isUploading && (
          <div className="status-message uploading">
            <p>Загрузка файла</p>
          </div>
        )}

        {uploadSuccess && (
          <div className="status-message success">
            <p>Файл успешно загружен</p>
          </div>
        )}

        {uploadError && (
          <div className="status-message error">
            <p>✗ {uploadError}</p>
          </div>
        )}

        {selectedFile && !uploadSuccess && (
          <button className="upload-button" onClick={handleUploadClick} disabled={isUploading}>
            {isUploading ? 'Загрузка...' : 'Загрузить'}
          </button>
        )}
      </div>
    </div>
  );
}
