using System.Windows;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;

namespace MediaCloudDesktop.Views
{
    public partial class WorkflowEditWindow : Window
    {
        private readonly ApiService _apiService;
        private readonly Workflow? _workflow;
        private readonly bool _isEditMode;

        public WorkflowEditWindow(ApiService apiService, Workflow? workflow)
        {
            InitializeComponent();
            _apiService = apiService;
            _workflow = workflow;
            _isEditMode = workflow != null;

            if (_isEditMode)
            {
                Title = "编辑工作流";
                NameTextBox.Text = workflow!.Name;
                DescriptionTextBox.Text = workflow.Description;
                TypeTextBox.Text = workflow.Type;
                ColorTextBox.Text = workflow.Color;
                IsActiveCheckBox.IsChecked = workflow.IsActive;
            }
            else
            {
                Title = "创建工作流";
            }
        }

        private async void SaveButton_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                SaveButton.IsEnabled = false;
                var name = NameTextBox.Text.Trim();
                var description = DescriptionTextBox.Text.Trim();
                var type = TypeTextBox.Text.Trim();
                var color = ColorTextBox.Text.Trim();
                var isActive = IsActiveCheckBox.IsChecked ?? true;

                if (string.IsNullOrEmpty(name))
                {
                    MessageBox.Show("请输入名称", "提示", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                if (_isEditMode)
                {
                    await _apiService.UpdateWorkflowAsync(_workflow!.Id, new { name, description, type, color, is_active = isActive });
                    MessageBox.Show("更新成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                }
                else
                {
                    await _apiService.CreateWorkflowAsync(name, description, type, color, isActive);
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
    }
}
