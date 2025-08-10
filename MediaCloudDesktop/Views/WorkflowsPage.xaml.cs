using System.Windows;
using System.Windows.Controls;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.Collections.ObjectModel;

namespace MediaCloudDesktop.Views
{
    public partial class WorkflowsPage : Page
    {
        private readonly ApiService _apiService;
        private readonly User _currentUser;
        private readonly ObservableCollection<Workflow> _workflows;
        private int _currentPage = 1;
        private int _pageSize = 20;
        private int _totalPages = 1;
        private string _searchKeyword = "";

        public WorkflowsPage(ApiService apiService, User currentUser)
        {
            InitializeComponent();
            _apiService = apiService;
            _currentUser = currentUser;
            _workflows = new ObservableCollection<Workflow>();
            WorkflowsDataGrid.ItemsSource = _workflows;
            
            LoadWorkflows();
        }

        private async void LoadWorkflows()
        {
            try
            {
                var response = await _apiService.GetWorkflowsAsync(_currentPage, _pageSize, _searchKeyword);
                
                _workflows.Clear();
                if (response.Data != null)
                {
                    foreach (var workflow in response.Data)
                    {
                        _workflows.Add(workflow);
                    }
                }

                if (response.Pagination != null)
                {
                    _totalPages = (int)Math.Ceiling((double)response.Pagination.Total / _pageSize);
                    PageInfoTextBlock.Text = $"第 {_currentPage} 页，共 {_totalPages} 页";
                }

                UpdatePaginationButtons();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载工作流失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private void UpdatePaginationButtons()
        {
            PrevPageButton.IsEnabled = _currentPage > 1;
            NextPageButton.IsEnabled = _currentPage < _totalPages;
        }

        private void CreateWorkflowButton_Click(object sender, RoutedEventArgs e)
        {
            var workflowWindow = new WorkflowEditWindow(_apiService, null);
            if (workflowWindow.ShowDialog() == true)
            {
                LoadWorkflows();
            }
        }

        private void SearchTextBox_TextChanged(object sender, TextChangedEventArgs e)
        {
            _searchKeyword = SearchTextBox.Text.Trim();
            _currentPage = 1;
            LoadWorkflows();
        }

        private void RefreshButton_Click(object sender, RoutedEventArgs e)
        {
            LoadWorkflows();
        }

        private void PrevPageButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentPage > 1)
            {
                _currentPage--;
                LoadWorkflows();
            }
        }

        private void NextPageButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentPage < _totalPages)
            {
                _currentPage++;
                LoadWorkflows();
            }
        }

        private void WorkflowsDataGrid_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            // 可以在这里添加选中工作流的处理逻辑
        }

        private void ViewWorkflow_Click(object sender, RoutedEventArgs e)
        {
            if (sender is Button button && button.DataContext is Workflow workflow)
            {
                var workflowWindow = new WorkflowDetailWindow(_apiService, workflow);
                workflowWindow.ShowDialog();
            }
        }

        private void EditWorkflow_Click(object sender, RoutedEventArgs e)
        {
            if (sender is Button button && button.DataContext is Workflow workflow)
            {
                // 检查权限：只能编辑自己创建的工作流或管理员
                if (_currentUser.Role != "admin" && workflow.CreatedBy != _currentUser.Id)
                {
                    MessageBox.Show("您没有权限编辑此工作流", "权限不足", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                var workflowWindow = new WorkflowEditWindow(_apiService, workflow);
                if (workflowWindow.ShowDialog() == true)
                {
                    LoadWorkflows();
                }
            }
        }

        private async void DeleteWorkflow_Click(object sender, RoutedEventArgs e)
        {
            if (sender is Button button && button.DataContext is Workflow workflow)
            {
                // 检查权限：只能删除自己创建的工作流或管理员
                if (_currentUser.Role != "admin" && workflow.CreatedBy != _currentUser.Id)
                {
                    MessageBox.Show("您没有权限删除此工作流", "权限不足", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }

                var result = MessageBox.Show($"确定要删除工作流 '{workflow.Name}' 吗？", "确认删除", 
                    MessageBoxButton.YesNo, MessageBoxImage.Question);
                
                if (result == MessageBoxResult.Yes)
                {
                    try
                    {
                        await _apiService.DeleteWorkflowAsync(workflow.Id);
                        MessageBox.Show("删除成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                        LoadWorkflows();
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
