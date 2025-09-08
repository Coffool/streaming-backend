def serialize_song(s):
    # s: (id, title, duration, audio_url, album_id)
    return {
        "id": s[0],
        "title": s[1],
        "duration": s[2],
        "audio_url": s[3],
    }


def serialize_album(a):
    return {
        "id": a[0],
        "title": a[1],
        "cover_url": a[2],
    }


def serialize_artist(ar):
    return {
        "id": ar[0],
        "name": ar[1],
        "profile_pic": ar[2],
    }
