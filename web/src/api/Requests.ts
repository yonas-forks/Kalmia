import { AxiosError, AxiosResponse, Method } from "axios";

import { FooterLabelLinks, LanderDetails, MoreLabelLinks } from "../types/doc";
import instance, { makeRequestWithCustomAuth } from "./AxiosInstance";

interface ErrorMessages {
  [key: number]: [string, string];
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface ApiResponse<T = any> {
  status: "success" | "error";
  code: number;
  data: T | null;
  message: string;
  path: string;
}

interface AuthCredentials {
  username: string;
  password: string;
}

export interface DocumentationPayload {
  id: number | null;
  name: string;
  version: string;
  url: string;
  organizationName: string;
  projectName: string;
  landerDetails: LanderDetails | string;
  baseURL: string;
  description: string;
  favicon: string;
  metaImage: string;
  navImage: string;
  navImageDark: string;
  customCSS: string;
  footerLabelLinks?: FooterLabelLinks[] | string;
  moreLabelLinks?: MoreLabelLinks[] | string;
  copyrightText: string;
  requireAuth?: boolean;
  gitUser?: string;
  gitRepo?: string;
  gitEmail?: string;
  gitPassword?: string;
  gitBranch?: string;

  // stored in bucket named
  bucketFavicon: string;
  bucketMetaImage: string;
  bucketNavImage: string;
  bucketNavImageDark: string;
}

interface CreateVersionPayload {
  originalDocId: number;
  version: string;
}

interface PageGroupPayload {
  name?: string;
  documentationId?: number;
  parentId?: number;
  id?: number;
  order?: number;
}

interface PagePayload {
  title: string;
  slug: string;
  content: string;
  id?: number;
  documentationId?: number;
  order?: number;
  pageGroupId?: number;
}

export interface OrderItem {
  id: number;
  order: number;
  documentationId?: number;
  parentId?: number | null;
  pageGroupId?: number;
  isPageGroup?: boolean;
}

interface ReorderBulkDataPayload {
  order: OrderItem[];
}

export interface UserPayload {
  username: string;
  email: string;
  password: string;
  permissions: string[];
  admin?: boolean;
}

export interface GitBookPayload {
  username: string;
  password: string;
  url: string;
}

export interface UpdateUserPayload {
  id: number;
  photo?: string;
  username?: string;
  email?: string;
  password?: string;
  permissions?: string[];
  admin?: boolean | number;
}

const ERROR_MESSAGES: ErrorMessages = {
  500: ["Internal Server Error", "/505"],
  404: ["Not Found", "/404"],
  401: ["Unauthorized", "/401"],
};

function getMessageAndPath(statusCode: number): [string, string] {
  return ERROR_MESSAGES[statusCode] || ["Unknown Error", "/error"];
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
async function makeRequest<T = any>(
  url: string,
  method: Method = "get",
  data: any = null, // eslint-disable-line @typescript-eslint/no-explicit-any
): Promise<ApiResponse<T>> {
  try {
    let response: AxiosResponse<T>;

    switch (method.toLowerCase()) {
      case "get":
        response = await instance.get<T>(url, { params: data });
        break;
      case "post":
        response = await instance.post<T>(url, data);
        break;
      case "put":
        response = await instance.put<T>(url, data);
        break;
      case "delete":
        response = await instance.delete<T>(url, { data });
        break;
      case "patch":
        response = await instance.patch<T>(url, data);
        break;
      default:
        throw new Error(`Unsupported HTTP method: ${method}`);
    }

    return {
      status: "success",
      code: response.status,
      data: response.data,
      message: "",
      path: "",
    };
  } catch (error) {
    const axiosError = error as AxiosError;
    const statusCode = axiosError?.response?.status || 500;
    const errorData = axiosError?.response?.data || null;
    const [message, path] = getMessageAndPath(statusCode);

    return {
      status: "error",
      code: statusCode,
      data: errorData as T | null,
      message,
      path,
    };
  }
}

export const createJWT = (data: AuthCredentials): Promise<ApiResponse> =>
  makeRequest("/kal-api/auth/jwt/create", "post", data);

export const refreshJWT = (token: string | null): Promise<ApiResponse> =>
  makeRequest("/kal-api/auth/jwt/refresh", "post", { token });

export const validateJWT = async (token: string): Promise<ApiResponse> => {
  try {
    const response: AxiosResponse = await makeRequestWithCustomAuth(
      "/kal-api/auth/jwt/validate",
      "post",
      token,
    );
    return {
      status: "success" as const,
      code: response.status,
      data: response.data,
      message: "",
      path: "",
    };
  } catch (error) {
    const axiosError = error as AxiosError;
    const statusCode = axiosError?.response?.status || 500;
    const errorData = axiosError?.response?.data || null;
    const [message, path] = getMessageAndPath(statusCode);
    return {
      status: "error" as const,
      code: statusCode,
      data: errorData,
      message,
      path,
    };
  }
};

export const signOut = (token: string | null): Promise<ApiResponse> =>
  makeRequest("/kal-api/auth/jwt/revoke", "post", { token });

export const getDocumentations = () =>
  makeRequest("/kal-api/docs/documentations");

export const getDocumentation = (id: number) =>
  makeRequest("/kal-api/docs/documentation", "post", { id });

export const createDocumentation = (data: DocumentationPayload) =>
  makeRequest("/kal-api/docs/documentation/create", "post", data);

export const updateDocumentation = (data: DocumentationPayload) =>
  makeRequest("/kal-api/docs/documentation/edit", "post", data);

export const deleteDocumentation = (id: number) =>
  makeRequest("/kal-api/docs/documentation/delete", "post", { id });

export const createDocumentationVersion = (data: CreateVersionPayload) =>
  makeRequest("/kal-api/docs/documentation/version", "post", data);

export const buildTrigger = () => makeRequest("/kal-api/health/last-trigger");
export const getPageGroups = () => makeRequest("/kal-api/docs/page-groups");

export const getPageGroup = (id: number) =>
  makeRequest("/kal-api/docs/page-group", "post", { id });

export const createPageGroup = (data: PageGroupPayload) =>
  makeRequest("/kal-api/docs/page-group/create", "post", data);

export const updatePageGroup = (data: PageGroupPayload) =>
  makeRequest("/kal-api/docs/page-group/edit", "post", data);

export const deletePageGroup = (id: number) =>
  makeRequest("/kal-api/docs/page-group/delete", "post", { id });

export const getPages = () => makeRequest("/kal-api/docs/pages");

export const getPage = (id: number) =>
  makeRequest("/kal-api/docs/page", "post", { id });

export const createPage = (data: PagePayload) =>
  makeRequest("/kal-api/docs/page/create", "post", data);

export const updatePage = (data: PagePayload) =>
  makeRequest("/kal-api/docs/page/edit", "post", data);

export const deletePage = (id: number) =>
  makeRequest("/kal-api/docs/page/delete", "post", { id });

export const commonReorderBulk = (data: ReorderBulkDataPayload) =>
  makeRequest("/kal-api/docs/documentation/reorder-bulk", "post", data);

export const getUsers = () => makeRequest("/kal-api/auth/users");

export const getUser = (id: number) =>
  makeRequest("/kal-api/auth/user", "post", {
    id: Number.parseInt(id.toString()),
  });

export const createUser = (data: UserPayload) =>
  makeRequest("/kal-api/auth/user/create", "post", data);

export const updateUser = (data: UpdateUserPayload) =>
  makeRequest("/kal-api/auth/user/edit", "post", data);

export const uploadFile = (data: FormData) =>
  makeRequest("/kal-api/auth/user/upload-file", "post", data);

export const uploadAssetsFile = (data: FormData) =>
  makeRequest("/kal-api/auth/user/assets/upload-file", "post", data);

export const deleteUser = (username: string) =>
  makeRequest("/kal-api/auth/user/delete", "post", { username });

export const getRootParentId = (docId: number) =>
  makeRequest(`/docs/documentation/root-parent-id?id=${docId}`);

export const importGitBook = (data: GitBookPayload) =>
  makeRequest("/kal-api/docs/import/gitbook", "post", data);

export const oAuthProviders = async (): Promise<string[]> => {
  const response = await makeRequest<string[]>("/kal-api/oauth/providers");
  if (response.status === "success") {
    return response.data as string[];
  } else {
    throw new Error(`Error fetching OAuth providers: ${response.message}`);
  }
};
