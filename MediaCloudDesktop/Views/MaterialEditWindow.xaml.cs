using System.Windows;
using System.Windows.Controls;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;

namespace MediaCloudDesktop.Views
{
    public partial class MaterialEditWindow : Window
    {
        private readonly ApiService _apiService;
        private readonly Material _material;
        private readonly List<Workflow> _workflows;

        public MaterialEditWindow(ApiService apiService, Material material, List<Workflow> workflows)
        {
            InitializeComponent();
            _apiService = apiService;
            _material = material;
            _workflows = workflows;

            NameTextBox.Text = material.OriginalFilename;
            IsStarredCheckBox.IsChecked = material.IsStarred;
            IsPublicCheckBox.IsChecked = material.IsPublic;

            foreach (var wf in _workflows)
            {
                var item = new ComboBoxItem { Content = wf.Name, Tag = wf.Id };
                if (material.WorkflowId.HasValue && wf.Id == material.WorkflowId.Value)
                {
                    item.IsSelected = true;
                }
                WorkflowComboBox.Items.Add(item);
            }

            // 加载全部标签，并选中已有标签
            _ = LoadTagsAsync();
        }

        private async Task LoadTagsAsync()
        {
            try
            {
                var tags = await _apiService.GetTagsAsync();
                TagsListBox.Items.Clear();
                var selectedIds = new HashSet<int>(_material.MaterialTags?.Select(mt => mt.Tag.Id) ?? Enumerable.Empty<int>());
                foreach (var tag in tags)
                {
                    var item = new ListBoxItem { Content = $"{tag.Name}", Tag = tag.Id };
                    if (selectedIds.Contains(tag.Id))
                    {
                        item.IsSelected = true;
                    }
                    TagsListBox.Items.Add(item);
                }
            }
            catch
            {
                // 忽略加载标签异常
            }
        }

        private async void SaveButton_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                var update = new Dictionary<string, object?>();
                if (!string.Equals(NameTextBox.Text.Trim(), _material.OriginalFilename))
                {
                    update["original_filename"] = NameTextBox.Text.Trim();
                }
                if (IsStarredCheckBox.IsChecked != _material.IsStarred)
                {
                    update["is_starred"] = IsStarredCheckBox.IsChecked;
                }
                if (IsPublicCheckBox.IsChecked != _material.IsPublic)
                {
                    update["is_public"] = IsPublicCheckBox.IsChecked;
                }

                // 处理工作流
                int? selectedWorkflowId = _material.WorkflowId;
                if (WorkflowComboBox.SelectedItem is ComboBoxItem sel)
                {
                    var tag = sel.Tag?.ToString();
                    if (tag == "__clear__")
                    {
                        selectedWorkflowId = null;
                    }
                    else if (int.TryParse(tag, out var wfId))
                    {
                        selectedWorkflowId = wfId;
                    }
                }
                if (selectedWorkflowId != _material.WorkflowId)
                {
                    update["workflow_id"] = selectedWorkflowId;
                }

                // 标签处理
                var selectedTagIds = new List<int>();
                foreach (var it in TagsListBox.SelectedItems.Cast<ListBoxItem>())
                {
                    if (it.Tag is int id)
                    {
                        selectedTagIds.Add(id);
                    }
                }
                if (_material.MaterialTags == null || !_material.MaterialTags.Select(mt => mt.Tag.Id).OrderBy(x=>x).SequenceEqual(selectedTagIds.OrderBy(x=>x)))
                {
                    update["tag_ids"] = selectedTagIds;
                }

                if (update.Count == 0)
                {
                    DialogResult = false;
                    return;
                }

                var updated = await _apiService.UpdateMaterialAsync(_material.Id, update);
                if (updated != null)
                {
                    DialogResult = true;
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"保存失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private void CancelButton_Click(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }
    }
}


