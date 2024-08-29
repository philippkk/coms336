ray: ray.cpp
	g++ ray.cpp -o ray -std=c++11
clean:
	rm -f *.o prog *~ core
