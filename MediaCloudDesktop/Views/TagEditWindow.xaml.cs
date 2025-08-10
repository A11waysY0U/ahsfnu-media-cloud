using System.Windows;
using System.Windows.Media;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;

namespace MediaCloudDesktop.Views
{
    public partial class TagEditWindow : Window
    {
        private readonly ApiService _apiService;
        private readonly Tag? _tag;
        private bool _isEditMode;

        public TagEditWindow(ApiService apiService, Tag? tag)
        {
            InitializeComponent();
            _apiService = apiService;
            _tag = tag;
            _isEditMode = tag != null;

            if (_isEditMode)
            {
                Title = "编辑标签";
                NameTextBox.Text = tag!.Name;
                ColorTextBox.Text = tag.Color;
                ColorPreview.Background = new SolidColorBrush((Color)ColorConverter.ConvertFromString(tag.Color));
            }
            else
            {
                Title = "创建标签";
            }
        }

        private void ColorPreview_MouseLeftButtonDown(object sender, System.Windows.Input.MouseButtonEventArgs e)
        {
            // 这里可以添加颜色选择器功能
            // 为了简化，我们使用预设的颜色
            var colors = new[] { "#409EFF", "#67C23A", "#E6A23C", "#F56C6C", "#909399", "#9C27B0", "#FF9800", "#795548" };
            var currentColor = ColorTextBox.Text;
            var currentIndex = Array.IndexOf(colors, currentColor);
            var nextIndex = (currentIndex + 1) % colors.Length;
            var nextColor = colors[nextIndex];
            
            ColorTextBox.Text = nextColor;
            ColorPreview.Background = new SolidColorBrush((Color)ColorConverter.ConvertFromString(nextColor));
        }

        private async void SaveButton_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                SaveButton.IsEnabled = false;
                var name = NameTextBox.Text.Trim();
                var color = ColorTextBox.Text.Trim();

                if (string.IsNullOrEmpty(name))
                {
                    MessageBox.Show("请输入标签名称", "提示", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                if (!IsValidColor(color))
                {
                    MessageBox.Show("请输入有效的颜色值（如：#409EFF）", "提示", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                if (_isEditMode)
                {
                    await _apiService.UpdateTagAsync(_tag!.Id, name, color);
                    MessageBox.Show("更新成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                }
                else
                {
                    await _apiService.CreateTagAsync(name, color);
                    MessageBox.Show("创建成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                }

                DialogResult = true;
                Close();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"操作失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
            finally
            {
                SaveButton.IsEnabled = true;
            }
        }

        private void CancelButton_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }

        private bool IsValidColor(string color)
        {
            try
            {
                ColorConverter.ConvertFromString(color);
                return true;
            }
            catch
            {
                return false;
            }
        }
    }
}
