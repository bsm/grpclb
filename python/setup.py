from setuptools import setup, find_packages

setup(
    name='grpclb',
    version='0.3.1',
    description='grpclb contains automatically generated code for calling a grpclb server',

    author='Black Square Media',
    url='https://github.com/bsm/grpclb',
    license='MIT',

    install_requires=['grpcio'],
    packages=find_packages(),
    zip_safe=False,

    classifiers=[
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
    ],
)
