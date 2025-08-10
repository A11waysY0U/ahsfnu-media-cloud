using System.Windows;
using MediaCloudDesktop.Models;
using MediaCloudDesktop.Services;

namespace MediaCloudDesktop.Views
{
    public partial class WorkflowDetailWindow : Window
    {
        public WorkflowDetailWindow(ApiService apiService, Workflow workflow)
        {
            InitializeComponent();
            DataContext = workflow;
        }

        private void CloseButton_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }
    }
}
