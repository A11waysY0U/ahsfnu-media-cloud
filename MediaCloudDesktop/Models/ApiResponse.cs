using System.Text.Json.Serialization;

namespace MediaCloudDesktop.Models
{
    public class ApiResponse<T>
    {
        [JsonPropertyName("data")]
        public T? Data { get; set; }

        [JsonPropertyName("pagination")]
        public Pagination? Pagination { get; set; }

        [JsonPropertyName("message")]
        public string? Message { get; set; }

        [JsonPropertyName("error")]
        public string? Error { get; set; }
    }

    public class Pagination
    {
        [JsonPropertyName("page")]
        public int Page { get; set; }

        [JsonPropertyName("page_size")]
        public int PageSize { get; set; }

        [JsonPropertyName("total")]
        public int Total { get; set; }
    }

    public class LoginRequest
    {
        [JsonPropertyName("username")]
        public string Username { get; set; } = string.Empty;

        [JsonPropertyName("password")]
        public string Password { get; set; } = string.Empty;

        [JsonPropertyName("auth_token")]
        public string AuthToken { get; set; } = string.Empty;
    }

    public class RegisterRequest
    {
        [JsonPropertyName("username")]
        public string Username { get; set; } = string.Empty;

        [JsonPropertyName("email")]
        public string Email { get; set; } = string.Empty;

        [JsonPropertyName("password")]
        public string Password { get; set; } = string.Empty;

        [JsonPropertyName("invite_code")]
        public string InviteCode { get; set; } = string.Empty;

        [JsonPropertyName("auth_token")]
        public string AuthToken { get; set; } = string.Empty;
    }

    public class CaptchaResponse
    {
        [JsonPropertyName("captcha_id")]
        public string CaptchaId { get; set; } = string.Empty;

        [JsonPropertyName("captcha_b64")]
        public string CaptchaB64 { get; set; } = string.Empty;

        [JsonPropertyName("auth_token")]
        public string AuthToken { get; set; } = string.Empty;
    }

    public class CaptchaVerifyRequest
    {
        [JsonPropertyName("captcha_id")]
        public string CaptchaId { get; set; } = string.Empty;

        [JsonPropertyName("captcha_code")]
        public string CaptchaCode { get; set; } = string.Empty;
    }

    public class LoginResponse
    {
        [JsonPropertyName("token")]
        public string Token { get; set; } = string.Empty;

        [JsonPropertyName("user")]
        public User User { get; set; } = new();
    }
}
