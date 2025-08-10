using System.Windows;
using System.Windows.Controls;
using Microsoft.Win32;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Windows.Data;
using System.Windows.Input;
using System.Linq;
using System.IO;
using System.Windows.Media.Imaging;
using MediaCloudDesktop.Models;

namespace MediaCloudDesktop.Views
{
    public partial class MaterialsPage : Page
    {
        private readonly ApiService _apiService;
        private readonly User _currentUser;
        private readonly ObservableCollection<Material> _materials;
        private bool _isInitialized = false;
        private int _currentPage = 1;
        private int _pageSize = 20;
        private int _totalPages = 1;
        private string _searchKeyword = "";
        private string _fileTypeFilter = "";
        private int? _workflowFilter = null;
        private List<int> _tagFilterIds = new();
        private List<Workflow> _workflows = new();
        private bool _useCardView = false;
        private Point _dragStartPoint;

        public MaterialsPage()
        {
            MessageBox.Show("请勿直接导航到 MaterialsPage.xaml，必须通过代码传递 ApiService 和 User 参数！", "导航错误", MessageBoxButton.OK, MessageBoxImage.Error);
            throw new InvalidOperationException("请使用带参数的构造函数 MaterialsPage(ApiService, User)");
        }

        public MaterialsPage(ApiService apiService, User currentUser)
        {
            _apiService = apiService;
            _currentUser = currentUser;
            _materials = new ObservableCollection<Material>();

            InitializeComponent();

            MaterialsDataGrid.ItemsSource = _materials;
            CardItemsControl.ItemsSource = _materials;

            _isInitialized = true;

            // 填充上传工作流下拉
            UploadWorkflowComboBox.Items.Clear();
            UploadWorkflowComboBox.Items.Add(new ComboBoxItem { Content = "不指定工作流", Tag = "__none__", IsSelected = true });

            LoadWorkflows();
            LoadMaterials();
        }

        private async void LoadWorkflows()
        {
            try
            {
                var workflowsResponse = await _apiService.GetWorkflowsAsync(1, 1000);
                _workflows = workflowsResponse.Data ?? new List<Workflow>();
                
                WorkflowComboBox.Items.Clear();
                WorkflowComboBox.Items.Add(new ComboBoxItem { Content = "全部工作流", IsSelected = true });
                
                foreach (var workflow in _workflows)
                {
                    WorkflowComboBox.Items.Add(new ComboBoxItem { Content = workflow.Name, Tag = workflow.Id });
                    UploadWorkflowComboBox.Items.Add(new ComboBoxItem { Content = workflow.Name, Tag = workflow.Id });
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载工作流失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private async void LoadMaterials()
        {
            ApiResponse<List<Material>> response = null;
            try
            {
                // 显示加载指示
                Mouse.OverrideCursor = Cursors.Wait;
                var tagsParam = _tagFilterIds.Count > 0 ? string.Join(",", _tagFilterIds) : null;
                response = await _apiService.GetMaterialsAsync(
                    _currentPage, 
                    _pageSize, 
                    _searchKeyword, 
                    _fileTypeFilter, 
                    _workflowFilter,
                    tagsParam);

                _materials.Clear();
                if (response != null && response.Data != null)
                {
                    foreach (var material in response.Data)
                    {
                        _materials.Add(material);
                    }
                }
                if (response != null && response.Pagination != null && PageInfoTextBlock != null)
                {
                    _totalPages = (int)Math.Ceiling((double)response.Pagination.Total / _pageSize);
                    PageInfoTextBlock.Text = $"第 {_currentPage} 页，共 {_totalPages} 页";
                }
                else
                {
                    Console.WriteLine($"response.Pagination 或 PageInfoTextBlock 为 null。response.Pagination: {response?.Pagination}, PageInfoTextBlock: {PageInfoTextBlock != null}");
                }

                UpdatePaginationButtons();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载素材失败: {ex.Message}\n{ex.StackTrace}\nresponse: " + (response == null ? "null" : response.ToString()) + "\nresponse.Data: " + (response?.Data == null ? "null" : response.Data.ToString()) + "\nresponse.Pagination: " + (response?.Pagination == null ? "null" : response.Pagination.ToString()), "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                Console.WriteLine($"[异常] 加载素材失败: {ex}\n{ex.StackTrace}\nresponse: " + (response == null ? "null" : response.ToString()) + "\nresponse.Data: " + (response?.Data == null ? "null" : response.Data.ToString()) + "\nresponse.Pagination: " + (response?.Pagination == null ? "null" : response.Pagination.ToString()));
            }
            finally
            {
                Mouse.OverrideCursor = null;
            }
        }

        private void UpdatePaginationButtons()
        {
            if (PrevPageButton != null && NextPageButton != null)
            {
                PrevPageButton.IsEnabled = _currentPage > 1;
                NextPageButton.IsEnabled = _currentPage < _totalPages;
            }
        }

        private async void UploadButton_Click(object sender, RoutedEventArgs e)
        {
            var openFileDialog = new OpenFileDialog
            {
                Title = "选择要上传的文件",
                Multiselect = false,
                Filter = "所有支持的文件|*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.mp4;*.avi;*.mov;*.wmv;*.mp3;*.wav;*.flac;*.aac;*.pdf;*.doc;*.docx;*.xls;*.xlsx;*.ppt;*.pptx|图片文件|*.jpg;*.jpeg;*.png;*.gif;*.bmp|视频文件|*.mp4;*.avi;*.mov;*.wmv|音频文件|*.mp3;*.wav;*.flac;*.aac|文档文件|*.pdf;*.doc;*.docx;*.xls;*.xlsx;*.ppt;*.pptx|所有文件|*.*"
            };

            if (openFileDialog.ShowDialog() == true)
            {
                try
                {
                    UploadButton.IsEnabled = false;
                    UploadButton.Content = "上传中 0%...";

                    int? wfId = null;
                    if (UploadWorkflowComboBox.SelectedItem is ComboBoxItem sel && sel.Tag != null && int.TryParse(sel.Tag.ToString(), out int id))
                    {
                        wfId = id;
                    }

                    var material = await _apiService.UploadMaterialAsync(openFileDialog.FileName, wfId, p =>
                    {
                        Dispatcher.Invoke(() => UploadButton.Content = $"上传中 {p}%...");
                    });
                    
                    MessageBox.Show("上传成功！", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                    LoadMaterials();
                }
                catch (Exception ex)
                {
                    MessageBox.Show($"上传失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                }
                finally
                {
                    UploadButton.IsEnabled = true;
                    UploadButton.Content = "上传素材";
                }
            }
        }

        private void MaterialsDropArea_DragOver(object sender, DragEventArgs e)
        {
            if (e.Data.GetDataPresent(DataFormats.FileDrop))
            {
                e.Effects = DragDropEffects.Copy;
            }
            else
            {
                e.Effects = DragDropEffects.None;
            }
            e.Handled = true;
        }

        private async void MaterialsDropArea_Drop(object sender, DragEventArgs e)
        {
            if (!e.Data.GetDataPresent(DataFormats.FileDrop)) return;
            string[] files = (string[])e.Data.GetData(DataFormats.FileDrop);
            if (files.Length == 0) return;

            try
            {
                UploadButton.IsEnabled = false;
                UploadButton.Content = "上传中 0%...";

                int? wfId = null;
                if (UploadWorkflowComboBox.SelectedItem is ComboBoxItem sel && sel.Tag != null && int.TryParse(sel.Tag.ToString(), out int id))
                {
                    wfId = id;
                }

                var file = files[0];
                await _apiService.UploadMaterialAsync(file, wfId, p =>
                {
                    Dispatcher.Invoke(() => UploadButton.Content = $"上传中 {p}%...");
                });
                MessageBox.Show("上传成功！", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                LoadMaterials();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"上传失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
            finally
            {
                UploadButton.IsEnabled = true;
                UploadButton.Content = "上传素材";
            }
        }

        private void SelectFolderButton_Click(object sender, RoutedEventArgs e)
        {
            var dlg = new Microsoft.Win32.OpenFileDialog
            {
                Title = "选择任意一个文件以定位文件夹",
                CheckFileExists = true,
                Multiselect = true
            };
            if (dlg.ShowDialog() == true)
            {
                var folder = System.IO.Path.GetDirectoryName(dlg.FileName) ?? string.Empty;
                LocalFolderPathText.Text = folder;
                try
                {
                    var files = Directory.GetFiles(folder)
                        .Where(p => IsSupportedFile(p))
                        .Select(p => CreateLocalFileItem(p))
                        .ToList();
                    LocalFilesListBox.ItemsSource = files;
                }
                catch (Exception ex)
                {
                    MessageBox.Show($"读取文件夹失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                }
            }
        }

        private static bool IsSupportedFile(string path)
        {
            var ext = System.IO.Path.GetExtension(path).ToLowerInvariant();
            string[] supported = new[] { ".jpg",".jpeg",".png",".gif",".bmp",
                ".mp4",".avi",".mov",".wmv",".mp3",".wav",".flac",".aac",
                ".pdf",".doc",".docx",".xls",".xlsx",".ppt",".pptx"};
            return supported.Contains(ext);
        }

        private void LocalFilesListBox_PreviewMouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            _dragStartPoint = e.GetPosition(null);
        }

        private void LocalFilesListBox_MouseMove(object sender, MouseEventArgs e)
        {
            if (e.LeftButton == MouseButtonState.Pressed)
            {
                var pos = e.GetPosition(null);
                if (Math.Abs(pos.X - _dragStartPoint.X) > SystemParameters.MinimumHorizontalDragDistance ||
                    Math.Abs(pos.Y - _dragStartPoint.Y) > SystemParameters.MinimumVerticalDragDistance)
                {
                    if (LocalFilesListBox.SelectedItems.Count > 0)
                    {
                        var files = LocalFilesListBox.SelectedItems.Cast<LocalFileItem>().Select(f => f.FullPath).ToArray();
                        var data = new DataObject(DataFormats.FileDrop, files);
                        DragDrop.DoDragDrop(LocalFilesListBox, data, DragDropEffects.Copy);
                    }
                }
            }
        }

        private LocalFileItem CreateLocalFileItem(string path)
        {
            var item = new LocalFileItem
            {
                Name = System.IO.Path.GetFileName(path),
                FullPath = path,
                IsImage = new[] { ".jpg",".jpeg",".png",".gif",".bmp" }.Contains(System.IO.Path.GetExtension(path).ToLowerInvariant())
            };
            if (item.IsImage)
            {
                try
                {
                    var bmp = new BitmapImage();
                    bmp.BeginInit();
                    bmp.CacheOption = BitmapCacheOption.OnLoad;
                    bmp.UriSource = new Uri(path);
                    bmp.DecodePixelWidth = 64;
                    bmp.EndInit();
                    bmp.Freeze();
                    item.Thumbnail = bmp;
                }
                catch { }
            }
            else
            {
                // 非图片可设置简单占位缩略图：此处留空由模板背景展示
            }
            return item;
        }

        private void SearchTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            if (!_isInitialized) return;
            _searchKeyword = SearchTextBox.Text.Trim();
            _currentPage = 1;
            LoadMaterials();
        }

        private void FileTypeComboBox_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            if (!_isInitialized) return;
            if (FileTypeComboBox.SelectedItem is ComboBoxItem selectedItem)
            {
                _fileTypeFilter = selectedItem.Content.ToString() switch
                {
                    "图片" => "image",
                    "视频" => "video",
                    "音频" => "audio",
                    "文档" => "document",
                    _ => ""
                };
                _currentPage = 1;
                LoadMaterials();
            }
        }

        private void WorkflowComboBox_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            if (!_isInitialized) return;
            if (WorkflowComboBox.SelectedItem is ComboBoxItem selectedItem)
            {
                if (selectedItem.Tag != null && int.TryParse(selectedItem.Tag.ToString(), out int workflowId))
                {
                    _workflowFilter = workflowId;
                }
                else
                {
                    _workflowFilter = null;
                }
                _currentPage = 1;
                LoadMaterials();
            }
        }

        

        private void CardViewToggle_Click(object sender, RoutedEventArgs e)
        {
            _useCardView = !_useCardView;
            GridViewContainer.Visibility = _useCardView ? Visibility.Collapsed : Visibility.Visible;
            CardViewContainer.Visibility = _useCardView ? Visibility.Visible : Visibility.Collapsed;
        }

        private async void TagFilterButton_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                var tags = await _apiService.GetTagsAsync();
                TagFilterListBox.Items.Clear();
                foreach (var tag in tags)
                {
                    var item = new ListBoxItem { Content = tag.Name, Tag = tag.Id };
                    if (_tagFilterIds.Contains(tag.Id)) item.IsSelected = true;
                    TagFilterListBox.Items.Add(item);
                }
                TagFilterPopup.IsOpen = true;
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载标签失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private void TagFilterConfirmButton_Click(object sender, RoutedEventArgs e)
        {
            _tagFilterIds = TagFilterListBox.SelectedItems.Cast<ListBoxItem>().Select(i => (int)i.Tag).ToList();
            TagFilterPopup.IsOpen = false;
            _currentPage = 1;
            LoadMaterials();
        }

        private void TagFilterClearButton_Click(object sender, RoutedEventArgs e)
        {
            _tagFilterIds.Clear();
            TagFilterPopup.IsOpen = false;
            _currentPage = 1;
            LoadMaterials();
        }

        private void RefreshButton_Click(object sender, RoutedEventArgs e)
        {
            LoadMaterials();
        }

        private void PrevPageButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentPage > 1)
            {
                _currentPage--;
                LoadMaterials();
            }
        }

        private void NextPageButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentPage < _totalPages)
            {
                _currentPage++;
                LoadMaterials();
            }
        }

        private void MaterialsDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            if (MaterialsDataGrid.SelectedItem is Material material)
            {
                // 预留选择变化
            }
            SelectedCountText.Text = $"已选择 {MaterialsDataGrid.SelectedItems.Count} 项";
        }

        private void EditMaterial_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as FrameworkElement)?.DataContext is Material material)
            {
                var dialog = new MaterialEditWindow(_apiService, material, _workflows);
                dialog.Owner = Window.GetWindow(this);
                var result = dialog.ShowDialog();
                if (result == true)
                {
                    LoadMaterials();
                }
            }
        }

        private void ViewMaterial_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as FrameworkElement)?.DataContext is Material material)
            {
                var dialog = new MaterialDetailWindow(material);
                dialog.Owner = Window.GetWindow(this);
                dialog.ShowDialog();
            }
        }

        private async void DeleteMaterial_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as FrameworkElement)?.DataContext is Material material)
            {
                if (MessageBox.Show($"确认删除素材 {material.OriginalFilename}?", "确认删除", MessageBoxButton.YesNo, MessageBoxImage.Warning) == MessageBoxResult.Yes)
                {
                    try
                    {
                        await _apiService.DeleteMaterialAsync(material.Id);
                        LoadMaterials();
                    }
                    catch (Exception ex)
                    {
                        MessageBox.Show($"删除失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                    }
                }
            }
        }

        private async void StarToggle_Click(object sender, RoutedEventArgs e)
        {
            if ((sender as FrameworkElement)?.DataContext is Material material)
            {
                try
                {
                    var update = new { is_starred = !material.IsStarred };
                    await _apiService.UpdateMaterialAsync(material.Id, update);
                    LoadMaterials();
                }
                catch (Exception ex)
                {
                    MessageBox.Show($"操作失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                }
            }
        }

        private List<Material> GetSelectedMaterials()
        {
            return MaterialsDataGrid.SelectedItems.Cast<Material>().ToList();
        }

        private async void BatchStarButton_Click(object sender, RoutedEventArgs e)
        {
            var items = GetSelectedMaterials();
            if (items.Count == 0) return;
            foreach (var m in items)
            {
                try { await _apiService.UpdateMaterialAsync(m.Id, new { is_starred = true }); } catch { }
            }
            LoadMaterials();
        }

        private async void BatchUnstarButton_Click(object sender, RoutedEventArgs e)
        {
            var items = GetSelectedMaterials();
            if (items.Count == 0) return;
            foreach (var m in items)
            {
                try { await _apiService.UpdateMaterialAsync(m.Id, new { is_starred = false }); } catch { }
            }
            LoadMaterials();
        }

        private async void BatchPublicButton_Click(object sender, RoutedEventArgs e)
        {
            var items = GetSelectedMaterials();
            if (items.Count == 0) return;
            foreach (var m in items)
            {
                try { await _apiService.UpdateMaterialAsync(m.Id, new { is_public = true }); } catch { }
            }
            LoadMaterials();
        }

        private async void BatchPrivateButton_Click(object sender, RoutedEventArgs e)
        {
            var items = GetSelectedMaterials();
            if (items.Count == 0) return;
            foreach (var m in items)
            {
                try { await _apiService.UpdateMaterialAsync(m.Id, new { is_public = false }); } catch { }
            }
            LoadMaterials();
        }

        private async void BatchDeleteButton_Click(object sender, RoutedEventArgs e)
        {
            var items = GetSelectedMaterials();
            if (items.Count == 0) return;
            if (MessageBox.Show($"确认删除所选 {items.Count} 项?", "确认删除", MessageBoxButton.YesNo, MessageBoxImage.Warning) != MessageBoxResult.Yes) return;
            foreach (var m in items)
            {
                try { await _apiService.DeleteMaterialAsync(m.Id); } catch { }
            }
            LoadMaterials();
        }
    }
}
