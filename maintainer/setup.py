from setuptools import setup, find_packages

setup(
    name='summa-relay',
    version='0.2.0',
    description=('Summa minimal relay'),
    author=["James Prestwich"],
    license="LGPLv3.0",
    install_requires=[
        'aiohttp',
        'riemann-ether',
        'riemann-keys',
        'mypy-extensions'],
    packages=find_packages(),
    package_data={'relay': ['py.typed']},
    package_dir={'relay': 'relay'},
    python_requires='>=3.6'
)
