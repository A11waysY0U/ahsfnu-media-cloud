using System.Windows.Media;

namespace MediaCloudDesktop.Models
{
    public class LocalFileItem
    {
        public string Name { get; set; } = string.Empty;
        public string FullPath { get; set; } = string.Empty;
        public bool IsImage { get; set; }
        public ImageSource? Thumbnail { get; set; }
    }
}


