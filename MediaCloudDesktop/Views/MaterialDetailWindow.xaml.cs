using System.Windows;
using System.Windows.Media.Imaging;
using MediaCloudDesktop.Models;

namespace MediaCloudDesktop.Views
{
    public partial class MaterialDetailWindow : Window
    {
        private readonly Material _material;

        public MaterialDetailWindow(Material material)
        {
            InitializeComponent();
            _material = material;

            NameText.Text = material.OriginalFilename;
            TypeText.Text = $"类型: {material.FileType} ({material.MimeType})";
            SizeText.Text = $"大小: {material.FileSizeFormatted}";
            ResolutionText.Text = material.Width.HasValue && material.Height.HasValue ? $"分辨率: {material.Width}x{material.Height}" : "分辨率: -";
            DurationText.Text = material.Duration.HasValue ? $"时长: {material.Duration}s" : "时长: -";
            UploaderText.Text = $"上传者: {material.Uploader?.Username ?? "-"}";
            UploadTimeText.Text = $"上传时间: {material.UploadTime:yyyy-MM-dd HH:mm}";

            if (!string.IsNullOrEmpty(material.ThumbnailPath))
            {
                try
                {
                    var uri = new System.Uri(material.ThumbnailPath, UriKind.RelativeOrAbsolute);
                    if (!uri.IsAbsoluteUri)
                    {
                        uri = new System.Uri(new System.Uri(material.FilePath).GetLeftPart(System.UriPartial.Authority).TrimEnd('/') + "/" + material.ThumbnailPath.TrimStart('/'));
                    }
                    ThumbnailImage.Source = new BitmapImage(uri);
                }
                catch {}
            }
        }

        private void CopyLink_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                Clipboard.SetText(_material.FilePath ?? string.Empty);
                MessageBox.Show("已复制文件链接", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
            }
            catch {}
        }

        private void Close_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }
    }
}


