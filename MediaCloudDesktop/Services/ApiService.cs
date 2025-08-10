using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Text.Json;
using MediaCloudDesktop.Models;
using System.IO;

namespace MediaCloudDesktop.Services
{
    public class ApiService
    {
        private readonly HttpClient _httpClient;
        private readonly string _baseUrl = "http://localhost:8080/api/v1";
        private string? _authToken;
        private string ApiOrigin => new Uri("http://localhost:8080/api/v1").GetLeftPart(UriPartial.Authority);

        private string MakeAbsoluteUrl(string? path, bool isThumbnail)
        {
            if (string.IsNullOrWhiteSpace(path)) return string.Empty;
            var s = path.Replace("\\", "/");
            // 将相对缩略图路径补齐到 /uploads/
            if (isThumbnail)
            {
                if (!s.StartsWith("/uploads/", StringComparison.OrdinalIgnoreCase))
                {
                    if (s.StartsWith("uploads/", StringComparison.OrdinalIgnoreCase))
                    {
                        s = "/" + s;
                    }
                    else
                    {
                        s = "/uploads/" + s.TrimStart('/');
                    }
                }
            }
            else
            {
                if (!s.StartsWith("/")) s = "/" + s;
            }
            return ApiOrigin.TrimEnd('/') + s;
        }

        public ApiService()
        {
            _httpClient = new HttpClient();
            _httpClient.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
            _httpClient.Timeout = TimeSpan.FromSeconds(30);
        }

        public void SetAuthToken(string token)
        {
            _authToken = token;
            _httpClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);
        }

        public void ClearAuthToken()
        {
            _authToken = null;
            _httpClient.DefaultRequestHeaders.Authorization = null;
        }

        // 认证相关API
        public async Task<CaptchaResponse> GetCaptchaAsync()
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/auth/captcha");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<CaptchaResponse>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new CaptchaResponse();
        }

        public async Task<string> VerifyCaptchaAsync(string captchaId, string captchaCode)
        {
            var request = new CaptchaVerifyRequest
            {
                CaptchaId = captchaId,
                CaptchaCode = captchaCode
            };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            HttpResponseMessage? response = null;
            string? responseContent = null;
            try
            {
                response = await _httpClient.PostAsync($"{_baseUrl}/auth/verify-captcha", content);
                responseContent = await response.Content.ReadAsStringAsync();
                response.EnsureSuccessStatusCode();

                // 直接解析auth_token字段
                using var doc = JsonDocument.Parse(responseContent);
                if (doc.RootElement.TryGetProperty("auth_token", out var tokenElement))
                {
                    return tokenElement.GetString() ?? string.Empty;
                }
                return string.Empty;
            }
            catch (Exception ex)
            {
                throw new Exception($"验证码校验失败: {ex.Message}\n响应内容: {responseContent}");
            }
        }

        public async Task<LoginResponse> LoginAsync(string username, string password, string authToken)
        {
            var request = new LoginRequest
            {
                Username = username,
                Password = password,
                AuthToken = authToken
            };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            HttpResponseMessage? response = null;
            string? responseContent = null;
            try
            {
                response = await _httpClient.PostAsync($"{_baseUrl}/auth/login", content);
                responseContent = await response.Content.ReadAsStringAsync();
                response.EnsureSuccessStatusCode();
                return JsonSerializer.Deserialize<LoginResponse>(responseContent, new JsonSerializerOptions
                {
                    PropertyNameCaseInsensitive = true
                }) ?? new LoginResponse();
            }
            catch (Exception ex)
            {
                throw new Exception($"登录请求失败: {ex.Message}\n响应内容: {responseContent}");
            }
        }

        public async Task<LoginResponse> RegisterAsync(string username, string email, string password, string inviteCode, string authToken)
        {
            var request = new RegisterRequest
            {
                Username = username,
                Email = email,
                Password = password,
                InviteCode = inviteCode,
                AuthToken = authToken
            };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PostAsync($"{_baseUrl}/auth/register", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<LoginResponse>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new LoginResponse();
        }

        // 素材相关API
        public async Task<Material> UploadMaterialAsync(string filePath, int? workflowId = null, Action<int>? onProgress = null)
        {
            using var form = new MultipartFormDataContent();
            using var fileStream = File.OpenRead(filePath);
            var progressContent = new ProgressableStreamContent(fileStream, (sent, total) =>
            {
                if (total > 0)
                {
                    var percent = (int)Math.Round(sent * 100.0 / total);
                    onProgress?.Invoke(percent);
                }
            });
            var fileContent = new StreamContent(fileStream);
            form.Add(progressContent, "file", Path.GetFileName(filePath));
            
            if (workflowId.HasValue)
            {
                form.Add(new StringContent(workflowId.Value.ToString()), "workflow_id");
            }

            var response = await _httpClient.PostAsync($"{_baseUrl}/materials", form);
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Material>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Material();
        }

        public async Task<ApiResponse<List<Material>>> GetMaterialsAsync(int page = 1, int pageSize = 20, string? keyword = null, string? fileType = null, int? workflowId = null, string? tags = null)
        {
            var queryParams = new List<string>
            {
                $"page={page}",
                $"page_size={pageSize}"
            };

            if (!string.IsNullOrEmpty(keyword))
                queryParams.Add($"keyword={Uri.EscapeDataString(keyword)}");
            if (!string.IsNullOrEmpty(fileType))
                queryParams.Add($"file_type={Uri.EscapeDataString(fileType)}");
            if (workflowId.HasValue)
                queryParams.Add($"workflow_id={workflowId.Value}");
            if (!string.IsNullOrEmpty(tags))
                queryParams.Add($"tags={Uri.EscapeDataString(tags)}");

            var url = $"{_baseUrl}/materials?{string.Join("&", queryParams)}";
            try
            {
                var response = await _httpClient.GetAsync(url);
                response.EnsureSuccessStatusCode();
                var content = await response.Content.ReadAsStringAsync();
                Console.WriteLine($"[GetMaterialsAsync] 原始响应内容: {content}");
                var result = JsonSerializer.Deserialize<ApiResponse<List<Material>>>(content, new JsonSerializerOptions
                {
                    PropertyNameCaseInsensitive = true
                });
                if (result == null)
                {
                    Console.WriteLine($"[GetMaterialsAsync] 反序列化结果为null, 原始内容: {content}");
                    return new ApiResponse<List<Material>>();
                }
                // 规范化 URL：确保缩略图/文件地址为绝对地址
                if (result.Data != null)
                {
                    foreach (var m in result.Data)
                    {
                        if (!string.IsNullOrWhiteSpace(m.ThumbnailPath) && !m.ThumbnailPath.StartsWith("http", StringComparison.OrdinalIgnoreCase))
                        {
                            m.ThumbnailPath = MakeAbsoluteUrl(m.ThumbnailPath, isThumbnail: true);
                        }
                        if (!string.IsNullOrWhiteSpace(m.FilePath) && !m.FilePath.StartsWith("http", StringComparison.OrdinalIgnoreCase))
                        {
                            m.FilePath = MakeAbsoluteUrl(m.FilePath, isThumbnail: false);
                        }
                    }
                }
                return result;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"[GetMaterialsAsync] 异常: {ex}");
                return new ApiResponse<List<Material>>();
            }
        }

        public async Task<Material> GetMaterialAsync(int id)
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/materials/{id}");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var mat = JsonSerializer.Deserialize<Material>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Material();
            if (!string.IsNullOrWhiteSpace(mat.ThumbnailPath) && !mat.ThumbnailPath.StartsWith("http", StringComparison.OrdinalIgnoreCase))
            {
                mat.ThumbnailPath = MakeAbsoluteUrl(mat.ThumbnailPath, isThumbnail: true);
            }
            if (!string.IsNullOrWhiteSpace(mat.FilePath) && !mat.FilePath.StartsWith("http", StringComparison.OrdinalIgnoreCase))
            {
                mat.FilePath = MakeAbsoluteUrl(mat.FilePath, isThumbnail: false);
            }
            return mat;
        }

        public async Task<Material> UpdateMaterialAsync(int id, object updateData)
        {
            var json = JsonSerializer.Serialize(updateData);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PutAsync($"{_baseUrl}/materials/{id}", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<ApiResponse<Material>>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });
            return result?.Data ?? new Material();
        }

        public async Task DeleteMaterialAsync(int id)
        {
            var response = await _httpClient.DeleteAsync($"{_baseUrl}/materials/{id}");
            response.EnsureSuccessStatusCode();
        }

        // 标签相关API
        public async Task<List<Tag>> GetTagsAsync()
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/tags");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<ApiResponse<List<Tag>>>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });
            return result?.Data ?? new List<Tag>();
        }

        public async Task<Tag> CreateTagAsync(string name, string color = "#409EFF")
        {
            var request = new { name, color };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PostAsync($"{_baseUrl}/tags", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Tag>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Tag();
        }

        public async Task<Tag> UpdateTagAsync(int id, string name, string color)
        {
            var request = new { name, color };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PutAsync($"{_baseUrl}/tags/{id}", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Tag>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Tag();
        }

        public async Task DeleteTagAsync(int id)
        {
            var response = await _httpClient.DeleteAsync($"{_baseUrl}/tags/{id}");
            response.EnsureSuccessStatusCode();
        }

        // 工作流相关API
        public async Task<ApiResponse<List<Workflow>>> GetWorkflowsAsync(int page = 1, int pageSize = 20, string? keyword = null)
        {
            var queryParams = new List<string>
            {
                $"page={page}",
                $"page_size={pageSize}"
            };

            if (!string.IsNullOrEmpty(keyword))
                queryParams.Add($"keyword={Uri.EscapeDataString(keyword)}");

            var url = $"{_baseUrl}/workflows?{string.Join("&", queryParams)}";
            var response = await _httpClient.GetAsync(url);
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<ApiResponse<List<Workflow>>>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new ApiResponse<List<Workflow>>();
        }

        public async Task<Workflow> CreateWorkflowAsync(string name, string description = "", string type = "custom", string color = "#409EFF", bool isActive = true, string config = "", List<int>? members = null)
        {
            var request = new { name, description, type, color, is_active = isActive, config, members };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PostAsync($"{_baseUrl}/workflows", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Workflow>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Workflow();
        }

        public async Task<Workflow> GetWorkflowAsync(int id)
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/workflows/{id}");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Workflow>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Workflow();
        }

        public async Task<Workflow> UpdateWorkflowAsync(int id, object updateData)
        {
            var json = JsonSerializer.Serialize(updateData);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PutAsync($"{_baseUrl}/workflows/{id}", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<Workflow>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new Workflow();
        }

        public async Task DeleteWorkflowAsync(int id)
        {
            var response = await _httpClient.DeleteAsync($"{_baseUrl}/workflows/{id}");
            response.EnsureSuccessStatusCode();
        }

        // 用户相关API
        public async Task<User> GetProfileAsync()
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/profile");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<User>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new User();
        }

        public async Task<List<User>> GetUsersAsync(int page = 1, int pageSize = 20, string? keyword = null, string? role = null)
        {
            var queryParams = new List<string>
            {
                $"page={page}",
                $"page_size={pageSize}"
            };

            if (!string.IsNullOrWhiteSpace(keyword))
            {
                queryParams.Add($"keyword={Uri.EscapeDataString(keyword)}");
            }

            if (!string.IsNullOrWhiteSpace(role))
            {
                queryParams.Add($"role={Uri.EscapeDataString(role)}");
            }

            var url = $"{_baseUrl}/users?{string.Join("&", queryParams)}";
            var response = await _httpClient.GetAsync(url);
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<ApiResponse<List<User>>>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });
            return result?.Data ?? new List<User>();
        }

        // 邀请码相关API
        public async Task<List<InviteCode>> GenerateInviteCodesAsync(int count)
        {
            var request = new { count };
            var json = JsonSerializer.Serialize(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _httpClient.PostAsync($"{_baseUrl}/invite_codes", content);
            response.EnsureSuccessStatusCode();
            var responseContent = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<ApiResponse<List<InviteCode>>>(responseContent, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });
            return result?.Data ?? new List<InviteCode>();
        }

        public async Task<ApiResponse<List<InviteCode>>> GetInviteCodesAsync(int page = 1, int pageSize = 20)
        {
            var url = $"{_baseUrl}/invite_codes?page={page}&page_size={pageSize}";
            var response = await _httpClient.GetAsync(url);
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            return JsonSerializer.Deserialize<ApiResponse<List<InviteCode>>>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            }) ?? new ApiResponse<List<InviteCode>>();
        }

        public async Task<InviteCodeStats> GetInviteCodeStatsAsync()
        {
            var response = await _httpClient.GetAsync($"{_baseUrl}/invite_codes/stats");
            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<ApiResponse<InviteCodeStats>>(content, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });
            return result?.Data ?? new InviteCodeStats();
        }
    }
}
