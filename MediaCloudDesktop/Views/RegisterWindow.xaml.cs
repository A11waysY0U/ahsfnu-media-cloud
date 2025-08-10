using System.Windows;
using System.Windows.Input;
using System.Windows.Media.Imaging;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;
using System.IO;

namespace MediaCloudDesktop.Views
{
    public partial class RegisterWindow : Window
    {
        private readonly ApiService _apiService;
        private string _currentCaptchaId = string.Empty;
        private string _currentAuthToken = string.Empty;

        public RegisterWindow(ApiService apiService)
        {
            InitializeComponent();
            _apiService = apiService;
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
            }
            catch (Exception ex)
            {
                StatusTextBlock.Text = $"加载验证码失败: {ex.Message}";
            }
        }

        private void CaptchaImage_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            LoadCaptcha();
        }

        private async void RegisterButton_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                RegisterButton.IsEnabled = false;
                StatusTextBlock.Text = "正在注册...";

                var username = UsernameTextBox.Text.Trim();
                var email = EmailTextBox.Text.Trim();
                var password = PasswordBox.Password;
                var confirmPassword = ConfirmPasswordBox.Password;
                var inviteCode = InviteCodeTextBox.Text.Trim();
                var captchaCode = CaptchaTextBox.Text.Trim();

                // 验证输入
                if (string.IsNullOrEmpty(username) || string.IsNullOrEmpty(email) || 
                    string.IsNullOrEmpty(password) || string.IsNullOrEmpty(confirmPassword) ||
                    string.IsNullOrEmpty(inviteCode) || string.IsNullOrEmpty(captchaCode))
                {
                    StatusTextBlock.Text = "请填写完整信息";
                    return;
                }

                if (password != confirmPassword)
                {
                    StatusTextBlock.Text = "两次输入的密码不一致";
                    return;
                }

                if (password.Length < 6)
                {
                    StatusTextBlock.Text = "密码长度至少6位";
                    return;
                }

                // 验证验证码
                var verifiedAuthToken = await _apiService.VerifyCaptchaAsync(_currentCaptchaId, captchaCode);
                if (string.IsNullOrEmpty(verifiedAuthToken))
                {
                    StatusTextBlock.Text = "验证码错误";
                    LoadCaptcha();
                    return;
                }

                // 注册
                var registerResponse = await _apiService.RegisterAsync(username, email, password, inviteCode, verifiedAuthToken);
                if (string.IsNullOrEmpty(registerResponse.Token))
                {
                    StatusTextBlock.Text = "注册失败，请检查邀请码是否正确";
                    LoadCaptcha();
                    return;
                }

                MessageBox.Show("注册成功！", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                Close();
            }
            catch (Exception ex)
            {
                StatusTextBlock.Text = $"注册失败: {ex.Message}";
                LoadCaptcha();
            }
            finally
            {
                RegisterButton.IsEnabled = true;
            }
        }

        private void CancelButton_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }
    }
}
