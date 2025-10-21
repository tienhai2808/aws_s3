"use client";

import { useState, useCallback } from "react";
import { useDropzone, FileWithPath } from "react-dropzone";
import type { FC } from "react";
import Image from "next/image";
import {
  Upload,
  X,
  CheckCircle,
  AlertCircle,
  Loader2,
  Image as ImageIcon,
} from "lucide-react";
import { uploadImage } from "@/lib/api";

interface FileWithPreview extends File {
  id: string; 
  preview: string;
  uploadStatus?: "idle" | "uploading" | "success" | "error";
}

const ImageUploader: FC = () => {
  const [files, setFiles] = useState<FileWithPreview[]>([]);
  const [isUploading, setIsUploading] = useState(false);

  const onDrop = useCallback((acceptedFiles: FileWithPath[]) => {
    const newFiles = acceptedFiles.map((file) =>
      Object.assign(file, {
        id: `file-${Math.random().toString(36).substring(2, 9)}`,
        preview: URL.createObjectURL(file),
        uploadStatus: "idle" as const,
      })
    );
    setFiles((prev) => [...prev, ...newFiles]);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "image/*": [".jpeg", ".jpg", ".png", ".gif", ".webp"],
    },
    multiple: true,
  });

  const removeFile = (id: string) => {
    setFiles((prev) => {
      const fileToRemove = prev.find((file) => file.id === id);
      if (fileToRemove) {
        URL.revokeObjectURL(fileToRemove.preview);
      }
      return prev.filter((file) => file.id !== id);
    });
  };

  const handleUpload = async () => {
    const filesToUpload = files.filter(
      (file) => file.uploadStatus === "idle" || file.uploadStatus === "error"
    );
    if (filesToUpload.length === 0) return;

    setIsUploading(true);

    const filesToUploadIds = new Set(filesToUpload.map((f) => f.id));

    setFiles((prevFiles) =>
      prevFiles.map((file) =>
        filesToUploadIds.has(file.id)
          ? { ...file, uploadStatus: "uploading" as const }
          : file
      )
    );

    const uploadPromises = filesToUpload.map((file) =>
      uploadImage(file)
        .then(() => ({ id: file.id, status: "success" as const }))
        .catch((error) => {
          console.error(`Upload LỖI cho file ${file.name}:`, error);
          return { id: file.id, status: "error" as const };
        })
    );

    const results = await Promise.all(uploadPromises);

    const resultsMap = new Map(results.map((r) => [r.id, r.status]));

    setFiles((currentFiles) =>
      currentFiles.map((file) => {
        const newStatus = resultsMap.get(file.id);
        return newStatus ? { ...file, uploadStatus: newStatus } : file;
      })
    );

    setIsUploading(false);
  };

  const hasFiles = files.length > 0;
  const hasUnuploadedFiles = files.some(
    (f) => f.uploadStatus === "idle" || f.uploadStatus === "error"
  );

  return (
    <div className="w-full max-w-4xl mx-auto p-6">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold text-gray-100 mb-2">Upload Ảnh</h1>
        <p className="text-gray-600">
          Tải ảnh lên S3 một cách nhanh chóng và dễ dàng
        </p>
      </div>
      <div
        {...getRootProps()}
        className={`relative border-2 border-dashed rounded-2xl p-12 transition-all duration-200 cursor-pointer ${
          isDragActive
            ? "border-blue-500 bg-blue-50 scale-[1.02]"
            : "border-gray-300 hover:border-gray-400 hover:bg-gray-50"
        }`}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center justify-center text-center">
          <div
            className={`rounded-full p-4 mb-4 transition-colors ${
              isDragActive ? "bg-blue-100" : "bg-gray-100"
            }`}
          >
            <ImageIcon
              className={`w-12 h-12 ${
                isDragActive ? "text-blue-500" : "text-gray-400"
              }`}
            />
          </div>
          <p className="text-lg font-semibold text-gray-700 mb-2">
            {isDragActive ? "Thả ảnh vào đây" : "Kéo thả ảnh hoặc nhấn để chọn"}
          </p>
          <p className="text-sm text-gray-500">Hỗ trợ: JPG, PNG, GIF, WEBP</p>
        </div>
      </div>
      {hasFiles && (
        <div className="mt-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">
              Ảnh đã chọn ({files.length})
            </h2>
            {hasUnuploadedFiles && (
              <button
                onClick={handleUpload}
                disabled={isUploading}
                className="flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium rounded-xl transition-colors duration-200 shadow-lg shadow-blue-500/30 hover:shadow-xl"
              >
                {isUploading ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    Đang tải...
                  </>
                ) : (
                  <>
                    <Upload className="w-5 h-5" />
                    Upload tất cả
                  </>
                )}
              </button>
            )}
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            {files.map((file) => (
              <div
                key={file.id}
                className="group relative aspect-square rounded-xl overflow-hidden bg-gray-100 shadow-md hover:shadow-xl transition-shadow duration-200"
              >
                <Image
                  src={file.preview}
                  alt={`Preview of ${file.name}`}
                  fill
                  className="object-cover"
                  unoptimized
                />

                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/0 to-black/0 opacity-0 group-hover:opacity-100 transition-opacity duration-200" />

                <button
                  onClick={() => removeFile(file.id)}
                  className="absolute top-2 right-2 p-1.5 bg-red-500 hover:bg-red-600 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-200 shadow-lg"
                  disabled={file.uploadStatus === "uploading"}
                >
                  <X className="w-4 h-4" />
                </button>

                <div className="absolute bottom-0 left-0 right-0 p-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  <p className="text-xs text-white font-medium truncate">
                    {file.name}
                  </p>
                </div>

                <div className="absolute top-2 left-2">
                  {file.uploadStatus === "uploading" && (
                    <div className="flex items-center gap-1 px-2 py-1 bg-blue-500 text-white rounded-full text-xs font-medium shadow-lg">
                      <Loader2 className="w-3 h-3 animate-spin" />
                      <span>Đang tải</span>
                    </div>
                  )}
                  {file.uploadStatus === "success" && (
                    <div className="flex items-center gap-1 px-2 py-1 bg-green-500 text-white rounded-full text-xs font-medium shadow-lg">
                      <CheckCircle className="w-3 h-3" />
                      <span>Thành công</span>
                    </div>
                  )}
                  {file.uploadStatus === "error" && (
                    <div className="flex items-center gap-1 px-2 py-1 bg-red-500 text-white rounded-full text-xs font-medium shadow-lg">
                      <AlertCircle className="w-3 h-3" />
                      <span>Lỗi</span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default ImageUploader;