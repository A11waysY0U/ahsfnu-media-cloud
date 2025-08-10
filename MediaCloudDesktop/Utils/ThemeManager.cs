using System.Windows;
using System.Windows.Media;

namespace MediaCloudDesktop.Utils
{
    public static class ThemeManager
    {
        public static void ApplyLightTheme(ResourceDictionary resources)
        {
            SetBrush(resources, "PrimaryBrush", ColorFromHex("#409EFF"));
            SetBrush(resources, "PrimaryHoverBrush", ColorFromHex("#337ECC"));
            SetBrush(resources, "PrimaryPressedBrush", ColorFromHex("#2B5AA0"));
            SetBrush(resources, "TextPrimaryBrush", ColorFromHex("#303133"));
            SetBrush(resources, "TextSecondaryBrush", ColorFromHex("#606266"));
            SetBrush(resources, "BorderBrushLight", ColorFromHex("#E4E7ED"));
            SetBrush(resources, "InputBorderBrush", ColorFromHex("#DCDFE6"));
            Application.Current.MainWindow.Background = Brushes.White;
        }

        public static void ApplyDarkTheme(ResourceDictionary resources)
        {
            SetBrush(resources, "PrimaryBrush", ColorFromHex("#5B8AFF"));
            SetBrush(resources, "PrimaryHoverBrush", ColorFromHex("#7AA4FF"));
            SetBrush(resources, "PrimaryPressedBrush", ColorFromHex("#3D63C6"));
            SetBrush(resources, "TextPrimaryBrush", ColorFromHex("#EAEAEA"));
            SetBrush(resources, "TextSecondaryBrush", ColorFromHex("#B8B8B8"));
            SetBrush(resources, "BorderBrushLight", ColorFromHex("#3C3F44"));
            SetBrush(resources, "InputBorderBrush", ColorFromHex("#4A4D52"));
            Application.Current.MainWindow.Background = new SolidColorBrush(ColorFromHex("#1F2125"));
        }

        private static void SetBrush(ResourceDictionary resources, string key, Color color)
        {
            // 不直接修改已有画刷（可能已被冻结），改为替换为新实例
            resources[key] = new SolidColorBrush(color);
        }

        private static Color ColorFromHex(string hex)
        {
            return (Color)ColorConverter.ConvertFromString(hex);
        }
    }
}


