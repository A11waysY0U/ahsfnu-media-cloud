using System.Text.Json.Serialization;

namespace MediaCloudDesktop.Models
{
    public class Tag
    {
        [JsonPropertyName("id")]
        public int Id { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("color")]
        public string Color { get; set; } = "#409EFF";

        [JsonPropertyName("created_by")]
        public int CreatedBy { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("creator")]
        public User? Creator { get; set; }

        [JsonPropertyName("material_tags")]
        public List<MaterialTag>? MaterialTags { get; set; }
    }
}
