

#include <stdlib.h>
#include <stddef.h>
#include <stdio.h>
#include <math.h>

#include <GL/gl.h>

#include "ogl-math.h"




void _vPrintMatrix (GLfloat mMatrix[16]) {
	
	int iX, iY;
	for (iY = 0; iY < 4; iY ++) {
		for (iX = 0; iX < 4; iX ++) {
			printf(" %.4f ", mMatrix[4 * iX + iY]);
		}
		printf("\n");
	}
	
}




void vIdentity (GLfloat mMatrix[16]) {
	
	mMatrix[ 0] = 1.0f; mMatrix[ 1] = 0.0f; mMatrix[ 2] = 0.0f; mMatrix[ 3] = 0.0f;
	mMatrix[ 4] = 0.0f; mMatrix[ 5] = 1.0f; mMatrix[ 6] = 0.0f; mMatrix[ 7] = 0.0f;
	mMatrix[ 8] = 0.0f; mMatrix[ 9] = 0.0f; mMatrix[10] = 1.0f; mMatrix[11] = 0.0f;
	mMatrix[12] = 0.0f; mMatrix[13] = 0.0f; mMatrix[14] = 0.0f; mMatrix[15] = 1.0f;
	
}




void vMultiply (GLfloat mMatrix[16], GLfloat mMatrixA[16], GLfloat mMatrixB[16]) {
	
	GLfloat nSum;
	int iX, iY, iI;
	for (iY = 0; iY < 4; iY ++) {
		for (iX = 0; iX < 4; iX ++) {
			nSum = 0;
			for (iI = 0; iI < 4; iI ++) {
				nSum += mMatrixA[4 * iI + iY] * mMatrixB[4 * iX + iI];
			}
			mMatrix[4 * iX + iY] = nSum;
		}
	}
	
}




void vGLMultiply (GLfloat mMatrix[16], GLfloat mMatrixA[16], GLfloat mMatrixB[16]) {
	
	glMatrixMode(GL_MODELVIEW_MATRIX);
	glLoadIdentity();
	glMultMatrixf(mMatrixA);
	glMultMatrixf(mMatrixB);
	glGetFloatv(GL_MODELVIEW_MATRIX, mMatrix);
	
}




void vTranslate (GLfloat mMatrix[16], GLfloat vTranslation[3]) {
	
	mMatrix[12] = vTranslation[0];
	mMatrix[13] = vTranslation[1];
	mMatrix[14] = vTranslation[2];
	
}




void _vGLTranslate (GLfloat mMatrix[16], GLfloat vTranslation[3]) {
	
	glMatrixMode(GL_MODELVIEW_MATRIX);
	glLoadMatrixf(mMatrix);
	glTranslatef(vTranslation[0], vTranslation[1], vTranslation[2]);
	glGetFloatv(GL_MODELVIEW_MATRIX, mMatrix);
	
}




void vGLRotate (GLfloat mMatrix[16], GLfloat nRotation, GLfloat vAxis[3]) {
	
	glMatrixMode(GL_MODELVIEW_MATRIX);
	glLoadMatrixf(mMatrix);
	glRotatef(nRotation, vAxis[0], vAxis[1], vAxis[2]);
	glGetFloatv(GL_MODELVIEW_MATRIX, mMatrix);
	
}




void vGLRotateH (GLfloat mMatrix[16], GLfloat nRotation) {
	
	GLfloat vAxis[3] = {0.0f, 1.0f, 0.0f};
	vGLRotate(mMatrix, nRotation, vAxis);
	
}




void vGLRotateV (GLfloat mMatrix[16], GLfloat nRotation) {
	
	GLfloat vAxis[3] = {1.0f, 0.0f, 0.0f};
	vGLRotate(mMatrix, nRotation, vAxis);
	
}





static void vFrustum (
	GLfloat mMatrix[16],
	GLfloat nLeft,
	GLfloat nRight,
	GLfloat nBottom,
	GLfloat nTop,
	GLfloat nNear,
	GLfloat nFar
) {
	
	
	GLfloat nA = (nRight + nLeft) / (nRight - nLeft);
	GLfloat nB = (nTop + nBottom) / (nTop - nBottom);
	GLfloat nC = - (nFar + nNear) / (nFar - nNear);
	GLfloat nD = - (2 * nFar * nNear) / (nFar - nNear);
	
	mMatrix[0] = (2 * nNear) / (nRight - nLeft); mMatrix[4] = 0; mMatrix[8] = nA; mMatrix[12] = 0;
	mMatrix[1] = 0; mMatrix[5] = (2 * nNear) / (nTop - nBottom); mMatrix[9] = nB; mMatrix[13] = 0;
	mMatrix[2] = 0; mMatrix[6] = 0; mMatrix[10] = nC; mMatrix[14] = nD;
	mMatrix[3] = 0; mMatrix[7] = 0; mMatrix[11] = -1; mMatrix[15] = 0;
	
}




void vPerspective (
	GLfloat mMatrix[16],
	GLfloat nFOVdeg,
	GLfloat nRatio,
	GLfloat nNear,
	GLfloat nFar
) {
	
	GLfloat nH = tan(nFOVdeg * M_PI / 360) * nNear;
	GLfloat nW = nH * nRatio;
	
	vFrustum(mMatrix, -nW, nW, -nH, nH, nNear, nFar);
	
}




void vMultiplyVectorWithMatrix (GLfloat* pVector, GLfloat* pMat) {
    
    float x = pMat[0]*pVector[0] + pMat[4]*pVector[1] + pMat[ 8]*pVector[2] + pMat[12]*pVector[3] ;
    float y = pMat[1]*pVector[0] + pMat[5]*pVector[1] + pMat[ 9]*pVector[2] + pMat[13]*pVector[3] ;
    float z = pMat[2]*pVector[0] + pMat[6]*pVector[1] + pMat[10]*pVector[2] + pMat[14]*pVector[3] ;
    float w = pMat[3]*pVector[0] + pMat[7]*pVector[1] + pMat[11]*pVector[2] + pMat[15]*pVector[3] ;

    pVector[0] = x;
    pVector[1] = y;
    pVector[2] = z;
    pVector[3] = w;
    
}




void vec_cross (GLfloat *out_result, GLfloat *u, GLfloat *v) {
    
    out_result[0] = u[1]*v[2] - u[2]*v[1];
    out_result[1] = u[2]*v[0] - u[0]*v[2];
    out_result[2] = u[0]*v[1] - u[1]*v[0];
    
}




void vec_normalize (GLfloat *inout_v) {
    
    GLfloat rlen = 1.0f/vec_length(inout_v);
    inout_v[0] *= rlen;
    inout_v[1] *= rlen;
    inout_v[2] *= rlen;
    
}




GLfloat vec_length (GLfloat *v) {
    
    return sqrtf(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]);
    
}



