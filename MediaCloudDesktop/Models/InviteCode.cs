using System.Text.Json.Serialization;

namespace MediaCloudDesktop.Models
{
    public class InviteCode
    {
        [JsonPropertyName("id")]
        public int Id { get; set; }

        [JsonPropertyName("code")]
        public string Code { get; set; } = string.Empty;

        [JsonPropertyName("status")]
        public int Status { get; set; }

        [JsonPropertyName("created_by")]
        public int CreatedBy { get; set; }

        [JsonPropertyName("used_by")]
        public int? UsedBy { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        public bool IsUsed => Status == 1;
        public bool IsExpired => Status == 2;
    }

    public class InviteCodeStats
    {
        [JsonPropertyName("total")]
        public int Total { get; set; }

        [JsonPropertyName("unused")]
        public int Unused { get; set; }

        [JsonPropertyName("used")]
        public int Used { get; set; }

        [JsonPropertyName("expired")]
        public int Expired { get; set; }
    }
}
