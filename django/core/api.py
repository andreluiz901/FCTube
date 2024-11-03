from core.models import Video
from core.serializers import VideoSerializer
from rest_framework.decorators import api_view
from rest_framework.response import Response
from django.db.models import Q

@api_view(['GET'])
def videos_list(request):
    q = request.query_params.get('q')
    if q:
        videos = Video.objects.filter(Q(title__icontains=q) | Q(description__icontains=q) | Q(tags__name__icontains=q))
    else:
        videos = Video.objects.all()
    query = videos.filter(is_published=True).order_by('-published_at').distinct()
    serializer = VideoSerializer(query, many=True)
    return Response(serializer.data)

@api_view(['GET'])
def videos_detail_by_slug(request, slug):
    video = Video.objects.get(slug=slug)
    serializer = VideoSerializer(video)
    return Response(serializer.data)

@api_view(['GET'])
def videos_detail_by_id(request, id):
    print('aaaaaa')
    video = Video.objects.get(id=id)
    serializer = VideoSerializer(video)
    return Response(serializer.data)

@api_view(['GET'])
def videos_list_recommended(request, id):
    video = Video.objects.get(id=id)
    tags = video.tags.all()
    videos = Video.objects.filter(tags__in=tags).exclude(id=id).distinct().order_by('-num_views')
    serializer = VideoSerializer(videos, many=True)
    return Response(serializer.data)

@api_view(['GET'])
def videos_get_likes(request, id):
    video = Video.objects.get(id=id)
    return Response({'likes': video.num_likes})

@api_view(['GET'])
def videos_get_views(request, id):
    video = Video.objects.get(id=id)
    return Response({'views': video.num_views})

@api_view(['POST'])
def videos_add_like(request, id):
    video = Video.objects.get(id=id)
    video.num_likes += 1
    video.save()
    return Response({'likes': video.num_likes})

@api_view(['POST'])
def videos_add_unlike(request, id):
    video = Video.objects.get(id=id)
    video.num_likes -= 1
    video.save()
    return Response({'likes': video.num_likes})

@api_view(['POST'])
def videos_register_view(request, id):
    video = Video.objects.get(id=id)
    video.num_views += 1
    video.save()
    return Response({'views': video.num_views})