using System.Windows.Controls;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.Collections.ObjectModel;
using System.Windows;

namespace MediaCloudDesktop.Views
{
    public partial class InviteCodesPage : Page
    {
        private readonly ApiService _apiService;
        private readonly User _currentUser;
        private readonly ObservableCollection<InviteCode> _inviteCodes;

        public InviteCodesPage(ApiService apiService, User currentUser)
        {
            InitializeComponent();
            _apiService = apiService;
            _currentUser = currentUser;
            _inviteCodes = new ObservableCollection<InviteCode>();
            InviteCodesDataGrid.ItemsSource = _inviteCodes;
            LoadInviteCodes();
        }

        private async void LoadInviteCodes()
        {
            try
            {
                var response = await _apiService.GetInviteCodesAsync();
                _inviteCodes.Clear();
                if (response.Data != null)
                {
                    foreach (var code in response.Data)
                    {
                        _inviteCodes.Add(code);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载邀请码失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private async void GenerateButton_Click(object sender, System.Windows.RoutedEventArgs e)
        {
            var input = Microsoft.VisualBasic.Interaction.InputBox("请输入要生成的邀请码数量（1-100）", "生成邀请码", "10");
            if (int.TryParse(input, out int count) && count > 0 && count <= 100)
            {
                try
                {
                    await _apiService.GenerateInviteCodesAsync(count);
                    MessageBox.Show("生成成功", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                    LoadInviteCodes();
                }
                catch (Exception ex)
                {
                    MessageBox.Show($"生成失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                }
            }
        }
    }
}
