from django import forms

MAX_VIDEO_CHUNK_SIZE = 1 * 1024 * 1024  # 1 MB

class VideoChunkUploadForm(forms.Form):
    chunk = forms.FileField(required=True)  # O arquivo chunk
    chunkIndex = forms.IntegerField(min_value=0, required=True)  # Índice do chunk (deve ser >= 0)

    def clean_chunk(self):
        chunk = self.cleaned_data.get('chunk')
        
        if chunk.size > MAX_VIDEO_CHUNK_SIZE:
            raise forms.ValidationError('O arquivo deve ser um vídeo no formato MP4.')

        return chunk
    
class VideoChunkFinishUploadForm(forms.Form):
    fileName = forms.CharField(max_length=255, required=True)  # Nome do arquivo
    totalChunks = forms.IntegerField(min_value=1, required=True)  # Total de chunks (deve ser >= 1)