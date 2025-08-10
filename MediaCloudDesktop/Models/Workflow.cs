using System.Text.Json.Serialization;

namespace MediaCloudDesktop.Models
{
    public class Workflow
    {
        [JsonPropertyName("id")]
        public int Id { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("type")]
        public string Type { get; set; } = "custom";

        [JsonPropertyName("color")]
        public string Color { get; set; } = "#409EFF";

        [JsonPropertyName("is_active")]
        public bool IsActive { get; set; } = true;

        [JsonPropertyName("config")]
        public string Config { get; set; } = string.Empty;

        [JsonPropertyName("status")]
        public string Status { get; set; } = "active";

        [JsonPropertyName("created_by")]
        public int CreatedBy { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("creator")]
        public User? Creator { get; set; }

        [JsonPropertyName("members")]
        public List<WorkflowMember>? Members { get; set; }
    }

    public class WorkflowMember
    {
        [JsonPropertyName("user_id")]
        public int UserId { get; set; }

        [JsonPropertyName("role")]
        public string Role { get; set; } = "member";

        [JsonPropertyName("user")]
        public User? User { get; set; }
    }
}
