using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media.Imaging;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using MediaCloudDesktop.Views;
using Microsoft.Win32;
using System.IO;
using System.Runtime.InteropServices;
using System;
using MediaCloudDesktop.Utils;

namespace MediaCloudDesktop
{
    public partial class MainWindow : Window
    {
        [DllImport("kernel32.dll")]
        public static extern bool AllocConsole();

        private readonly ApiService _apiService;
        private User? _currentUser;
        private string _currentAuthToken = string.Empty;
        private string _currentCaptchaId = string.Empty;

        public MainWindow()
        {
#if DEBUG
            AllocConsole();
#endif
            InitializeComponent();
            _apiService = new ApiService();
            LoadCaptcha();
        }

        private async void LoadCaptcha()
        {
            try
            {
                var captcha = await _apiService.GetCaptchaAsync();
                _currentCaptchaId = captcha.CaptchaId;
                _currentAuthToken = captcha.AuthToken;

                // 将base64图片转换为BitmapImage
                var imageBytes = Convert.FromBase64String(captcha.CaptchaB64.Replace("data:image/png;base64,", ""));
                var bitmap = new BitmapImage();
                bitmap.BeginInit();
                bitmap.StreamSource = new MemoryStream(imageBytes);
                bitmap.EndInit();
                CaptchaImage.Source = bitmap;
                Console.WriteLine($"[验证码] 获取成功 CaptchaId={_currentCaptchaId}");
            }
            catch (Exception ex)
            {
                StatusTextBlock.Text = $"加载验证码失败: {ex.Message}";
                Console.WriteLine($"[验证码] 获取失败: {ex}");
            }
        }

        private void CaptchaImage_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            LoadCaptcha();
        }

        private async void LoginButton_Click(object sender, RoutedEventArgs e)
        {
            MessageBox.Show("LoginButton_Click 触发");
            Console.WriteLine("[登录] LoginButton_Click 触发");
            try
            {
                LoginButton.IsEnabled = false;
                StatusTextBlock.Text = "正在登录...";

                var username = UsernameTextBox.Text.Trim();
                var password = PasswordBox.Password;
                var captchaCode = CaptchaTextBox.Text.Trim();

                if (string.IsNullOrEmpty(username) || string.IsNullOrEmpty(password) || string.IsNullOrEmpty(captchaCode))
                {
                    MessageBox.Show("输入不完整，return");
                    StatusTextBlock.Text = "请填写完整信息";
                    Console.WriteLine("[登录] 输入不完整，return");
                    return;
                }

                MessageBox.Show("输入完整，准备验证验证码");
                Console.WriteLine($"[登录] 输入完整，准备验证验证码 CaptchaId={_currentCaptchaId} Code={captchaCode}");
                // 验证验证码
                var verifiedAuthToken = await _apiService.VerifyCaptchaAsync(_currentCaptchaId, captchaCode);
                MessageBox.Show($"验证码验证结果: '{verifiedAuthToken}'");
                Console.WriteLine($"[登录] 验证码验证结果: '{verifiedAuthToken}'");
                if (string.IsNullOrEmpty(verifiedAuthToken))
                {
                    MessageBox.Show("验证码错误，return");
                    StatusTextBlock.Text = "验证码错误";
                    LoadCaptcha();
                    Console.WriteLine("[登录] 验证码错误，return");
                    return;
                }

                MessageBox.Show("验证码正确，准备登录请求");
                Console.WriteLine("[登录] 验证码正确，准备登录请求");
                // 登录
                var loginResponse = await _apiService.LoginAsync(username, password, verifiedAuthToken);
                MessageBox.Show($"登录响应Token: '{loginResponse.Token}'");
                Console.WriteLine($"[登录] 登录响应Token: '{loginResponse.Token}'");
                if (string.IsNullOrEmpty(loginResponse.Token))
                {
                    MessageBox.Show("登录失败，Token为空，return");
                    StatusTextBlock.Text = "登录失败，请检查用户名和密码";
                    LoadCaptcha();
                    Console.WriteLine("[登录] 登录失败，Token为空，return");
                    return;
                }

                MessageBox.Show("登录成功，切换主界面");
                Console.WriteLine("[登录] 登录成功，切换主界面");
                _currentUser = loginResponse.User;
                _apiService.SetAuthToken(loginResponse.Token);

                // 切换到主界面
                LoginGrid.Visibility = Visibility.Collapsed;
                MainGrid.Visibility = Visibility.Visible;
                UserInfoTextBlock.Text = $"欢迎，{_currentUser.Username}";

                // 根据用户角色显示管理功能
                if (_currentUser.Role == "admin")
                {
                    UsersButton.Visibility = Visibility.Visible;
                    InviteCodesButton.Visibility = Visibility.Visible;
                }

                // 默认显示素材管理页面
                MainFrame.Navigate(new MaterialsPage(_apiService, _currentUser));

                StatusTextBlock.Text = "";
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.ToString(), "登录异常", MessageBoxButton.OK, MessageBoxImage.Error);
                StatusTextBlock.Text = $"登录失败: {ex.Message}";
                LoadCaptcha();
                Console.WriteLine($"[登录] 异常: {ex}");
            }
            finally
            {
                LoginButton.IsEnabled = true;
            }
        }

        private void RegisterButton_Click(object sender, RoutedEventArgs e)
        {
            // 切换到注册控件，内容放入 AuthContent
            var registerControl = new Views.RegisterControl(_apiService);
            registerControl.RequestBackToLogin += (_, __) =>
            {
                AuthContent.Content = LoginPanel;
                LoginPanel.Visibility = Visibility.Visible;
                BrandPanel.Visibility = Visibility.Visible;
            };
            LoginPanel.Visibility = Visibility.Collapsed;
            AuthContent.Content = registerControl;
            BrandPanel.Visibility = Visibility.Collapsed;
        }

        private void LogoutButton_Click(object sender, RoutedEventArgs e)
        {
            _apiService.ClearAuthToken();
            _currentUser = null;
            _currentAuthToken = string.Empty;
            
            LoginGrid.Visibility = Visibility.Visible;
            MainGrid.Visibility = Visibility.Collapsed;
            
            UsernameTextBox.Text = "";
            PasswordBox.Password = "";
            CaptchaTextBox.Text = "";
            StatusTextBlock.Text = "";
            
            LoadCaptcha();
        }

        private void MaterialsButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentUser != null)
            {
                MainFrame.Navigate(new MaterialsPage(_apiService, _currentUser));
            }
        }

        private void TagsButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentUser != null)
            {
                MainFrame.Navigate(new TagsPage(_apiService, _currentUser));
            }
        }

        private void WorkflowsButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentUser != null)
            {
                MainFrame.Navigate(new WorkflowsPage(_apiService, _currentUser));
            }
        }

        private void UsersButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentUser != null && _currentUser.Role == "admin")
            {
                MainFrame.Navigate(new UsersPage(_apiService, _currentUser));
            }
        }

        private void InviteCodesButton_Click(object sender, RoutedEventArgs e)
        {
            if (_currentUser != null && _currentUser.Role == "admin")
            {
                MainFrame.Navigate(new InviteCodesPage(_apiService, _currentUser));
            }
        }

        private void ThemeToggle_Click(object sender, RoutedEventArgs e)
        {
            if (ThemeToggle.IsChecked == true)
            {
                ThemeManager.ApplyDarkTheme(Application.Current.Resources);
            }
            else
            {
                ThemeManager.ApplyLightTheme(Application.Current.Resources);
            }
        }
    }
}