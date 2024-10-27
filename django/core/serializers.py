from core.models import Video
from rest_framework import serializers, settings


class VideoSerializer(serializers.ModelSerializer):
    id = serializers.IntegerField()
    title = serializers.CharField()
    description = serializers.CharField()
    slug = serializers.CharField()
    published_at = serializers.DateTimeField()
    views = serializers.IntegerField(source='num_views')
    likes = serializers.IntegerField(source='num_likes')
    tags = serializers.PrimaryKeyRelatedField(many=True, read_only=True)
    thumbnail = serializers.SerializerMethodField()
    video_url = serializers.SerializerMethodField()


    def get_thumbnail(self, obj):
        assets_url = 'http://localhost:9000/media/uploads' # settings.ASSETS_URL
        return f'{assets_url}/{obj.thumbnail}'
    
    def get_video_url(self, obj):
        assets_url = 'http://localhost:9000/media/uploads' # settings.ASSETS_URL
        return f'{assets_url}{obj.video_media.video_path}'
    
    class Meta:
        model = Video
        fields = ['id','title', 'description', 'slug', 'published_at', 'views', 'likes', 'tags', 'thumbnail', 'video_url']
