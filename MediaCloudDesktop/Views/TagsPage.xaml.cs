using System.Windows;
using System.Windows.Controls;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.Collections.ObjectModel;

namespace MediaCloudDesktop.Views
{
    public partial class TagsPage : Page
    {
        private readonly ApiService _apiService;
        private readonly User _currentUser;
        private readonly ObservableCollection<Tag> _tags;

        public TagsPage(ApiService apiService, User currentUser)
        {
            InitializeComponent();
            _apiService = apiService;
            _currentUser = currentUser;
            _tags = new ObservableCollection<Tag>();
            TagsDataGrid.ItemsSource = _tags;
            
            LoadTags();
        }

        private async void LoadTags()
        {
            try
            {
                var tags = await _apiService.GetTagsAsync();
                _tags.Clear();
                foreach (var tag in tags)
                {
                    _tags.Add(tag);
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载标签失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private void CreateTagButton_Click(object sender, RoutedEventArgs e)
        {
            var tagWindow = new TagEditWindow(_apiService, null);
            if (tagWindow.ShowDialog() == true)
            {
                LoadTags();
            }
        }

        private void TagsDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            // 可以在这里添加选中标签的处理逻辑
        }

        private void EditTag_Click(object sender, RoutedEventArgs e)
        {
            if (sender is Button button && button.DataContext is Tag tag)
            {
                // 检查权限：只能编辑自己创建的标签或管理员
                if (_currentUser.Role != "admin" && tag.CreatedBy != _currentUser.Id)
                {
                    MessageBox.Show("您没有权限编辑此标签", "权限不足", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                var tagWindow = new TagEditWindow(_apiService, tag);
                if (tagWindow.ShowDialog() == true)
                {
                    LoadTags();
                }
            }
        }

        private async void DeleteTag_Click(object sender, RoutedEventArgs e)
        {
            if (sender is Button button && button.DataContext is Tag tag)
            {
                // 检查权限：只能删除自己创建的标签或管理员
                if (_currentUser.Role != "admin" && tag.CreatedBy != _currentUser.Id)
                {
                    MessageBox.Show("您没有权限删除此标签", "权限不足", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                var result = MessageBox.Show($"确定要删除标签 '{tag.Name}' 吗？", "确认删除", 
                    MessageBoxButton.YesNo, MessageBoxImage.Question);
                
                if (result == MessageBoxResult.Yes)
                {
                    try
                    {
                        await _apiService.DeleteTagAsync(tag.Id);
                        MessageBox.Show("删除成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                        LoadTags();
                    }
                    catch (Exception ex)
                    {
                        MessageBox.Show($"删除失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                    }
                }
            }
        }
    }
}
