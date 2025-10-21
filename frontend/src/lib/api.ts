import axios, { AxiosInstance } from "axios";

interface PresignedUrlResponse {
  upload_url: string;
}

interface PresignedUrlRequest {
  file_name: string;
  content_type: string;
}

export interface ApiResponse<T = unknown> {
  message: string;
  data: T;
}

const axiosInstance: AxiosInstance = axios.create({
  baseURL: "http://localhost:5000",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
  timeout: 10000,
});

export async function getPresignedUrl(
  data: PresignedUrlRequest
): Promise<ApiResponse<PresignedUrlResponse>> {
  try {
    const response = await axiosInstance.post<ApiResponse<PresignedUrlResponse>>(
      "/files/generate-presigned",
      data
    );
    return response.data;
  } catch (error) {
    console.error("Error getting presigned URL:", error);
    throw error;
  }
}

export async function uploadToS3(uploadURL: string, file: File): Promise<void> {
  try {
    await axios.put(uploadURL, file, {
      headers: {
        "Content-Type": file.type,
      },
    });
  } catch (error) {
    console.error("Error uploading to S3:", error);
    throw error;
  }
}

export async function uploadImage(file: File): Promise<void> {
  const data: PresignedUrlRequest = {
    file_name: file.name,
    content_type: file.type,
  }
  const resp: ApiResponse<PresignedUrlResponse> = await getPresignedUrl(data);

  await uploadToS3(resp.data.upload_url, file);
}
