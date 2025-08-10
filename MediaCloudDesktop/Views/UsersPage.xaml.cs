using System.Windows.Controls;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.Collections.ObjectModel;
using System.Windows;

namespace MediaCloudDesktop.Views
{
    public partial class UsersPage : Page
    {
        private readonly ApiService _apiService;
        private readonly User _currentUser;
        private readonly ObservableCollection<User> _users;

        public UsersPage(ApiService apiService, User currentUser)
        {
            InitializeComponent();
            _apiService = apiService;
            _currentUser = currentUser;
            _users = new ObservableCollection<User>();
            UsersDataGrid.ItemsSource = _users;
            LoadUsers();
        }

        private async void LoadUsers()
        {
            try
            {
                var users = await _apiService.GetUsersAsync();
                _users.Clear();
                foreach (var user in users)
                {
                    _users.Add(user);
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"加载用户失败: {ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }
    }
}
