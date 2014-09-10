
eins

///#///

void vIdentity (GLfloat mMatrix[16]) ;
void vMultiply (GLfloat mMatrix[16], GLfloat mMatrixA[16], GLfloat mMatrixB[16]) ;
void vGLMultiply (GLfloat mMatrix[16], GLfloat mMatrixA[16], GLfloat mMatrixB[16]) ;
void vTranslate (GLfloat mMatrix[16], GLfloat vTranslation[3]) ;
void vGLRotate (GLfloat mMatrix[16], GLfloat nRotation, GLfloat vAxis[3]) ;
void vGLRotateH (GLfloat mMatrix[16], GLfloat nRotation) ;
void vGLRotateV (GLfloat mMatrix[16], GLfloat nRotation) ;
static void vFrustum (
	GLfloat mMatrix[16],
	GLfloat nLeft,
	GLfloat nRight,
	GLfloat nBottom,
	GLfloat nTop,
	GLfloat nNear,
	GLfloat nFar
) ;
void vPerspective (
	GLfloat mMatrix[16],
	GLfloat nFOVdeg,
	GLfloat nRatio,
	GLfloat nNear,
	GLfloat nFar
) ;
void vMultiplyVectorWithMatrix (GLfloat* pVector, GLfloat* pMat) ;
void vec_cross (GLfloat *out_result, GLfloat *u, GLfloat *v) ;
void vec_normalize (GLfloat *inout_v) ;
GLfloat vec_length (GLfloat *v) ;

///#///

drei
