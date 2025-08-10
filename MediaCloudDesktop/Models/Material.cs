using System.Text.Json.Serialization;

namespace MediaCloudDesktop.Models
{
    public class Material
    {
        [JsonPropertyName("id")]
        public int Id { get; set; }

        [JsonPropertyName("filename")]
        public string Filename { get; set; } = string.Empty;

        [JsonPropertyName("original_filename")]
        public string OriginalFilename { get; set; } = string.Empty;

        [JsonPropertyName("file_path")]
        public string FilePath { get; set; } = string.Empty;

        [JsonPropertyName("file_size")]
        public long FileSize { get; set; }

        [JsonPropertyName("file_type")]
        public string FileType { get; set; } = string.Empty;

        [JsonPropertyName("mime_type")]
        public string MimeType { get; set; } = string.Empty;

        [JsonPropertyName("width")]
        public int? Width { get; set; }

        [JsonPropertyName("height")]
        public int? Height { get; set; }

        [JsonPropertyName("duration")]
        public int? Duration { get; set; }

        [JsonPropertyName("uploaded_by")]
        public int UploadedBy { get; set; }

        [JsonPropertyName("upload_time")]
        public DateTime UploadTime { get; set; }

        [JsonPropertyName("is_starred")]
        public bool IsStarred { get; set; }

        [JsonPropertyName("is_public")]
        public bool IsPublic { get; set; }

        [JsonPropertyName("workflow_id")]
        public int? WorkflowId { get; set; }

        [JsonPropertyName("thumbnail_path")]
        public string? ThumbnailPath { get; set; }

        [JsonPropertyName("uploader")]
        public User? Uploader { get; set; }

        [JsonPropertyName("material_tags")]
        public List<MaterialTag>? MaterialTags { get; set; }

        public string FileSizeFormatted => FormatFileSize(FileSize);

        private static string FormatFileSize(long bytes)
        {
            string[] sizes = { "B", "KB", "MB", "GB", "TB" };
            double len = bytes;
            int order = 0;
            while (len >= 1024 && order < sizes.Length - 1)
            {
                order++;
                len = len / 1024;
            }
            return $"{len:0.##} {sizes[order]}";
        }
    }

    public class MaterialTag
    {
        [JsonPropertyName("tag")]
        public Tag Tag { get; set; } = new();
    }
}
