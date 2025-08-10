using System.Net.Http;
using System.Net;
using System;
using System.IO;
using System.Threading.Tasks;

namespace MediaCloudDesktop.Services
{
    public class ProgressableStreamContent : HttpContent
    {
        private const int DefaultBufferSize = 81920;
        private readonly Stream _source;
        private readonly int _bufferSize;
        private readonly Action<long, long>? _progress;

        public ProgressableStreamContent(Stream source, Action<long, long>? progress = null, int bufferSize = DefaultBufferSize)
        {
            _source = source;
            _bufferSize = bufferSize;
            _progress = progress;
        }

        protected override async Task SerializeToStreamAsync(Stream stream, TransportContext? context)
        {
            var buffer = new byte[_bufferSize];
            long totalBytesRead = 0;
            long? length = null;
            try
            {
                length = _source.CanSeek ? _source.Length : null;
            }
            catch { }

            int bytesRead;
            while ((bytesRead = await _source.ReadAsync(buffer, 0, buffer.Length)) > 0)
            {
                await stream.WriteAsync(buffer.AsMemory(0, bytesRead));
                totalBytesRead += bytesRead;
                _progress?.Invoke(totalBytesRead, length ?? -1);
            }
        }

        protected override bool TryComputeLength(out long length)
        {
            if (_source.CanSeek)
            {
                length = _source.Length;
                return true;
            }
            length = -1;
            return false;
        }
    }
}


