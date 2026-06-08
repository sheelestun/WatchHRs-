import cv2
import numpy as np
import torch
from facenet_pytorch import InceptionResnetV1

_recognizer = InceptionResnetV1(pretrained="vggface2").eval()


def get_face_embedding(face_image: np.ndarray) -> np.ndarray:
    face_rgb = cv2.cvtColor(face_image, cv2.COLOR_BGR2RGB)
    face_resized = cv2.resize(face_rgb, (160, 160))
    face_float = face_resized.astype(np.float32) / 255.0
    face_normalized = (face_float - 0.5) / 0.5
    face_tensor = torch.from_numpy(face_normalized.transpose(2, 0, 1)).unsqueeze(0)

    with torch.no_grad():
        embedding = _recognizer(face_tensor)

    return embedding.squeeze().numpy()
